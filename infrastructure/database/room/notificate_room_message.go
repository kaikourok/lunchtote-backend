package room

import (
	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"golang.org/x/sync/errgroup"
)

func (db *RoomRepository) NotificateRoomMessage(messageId int) (dto *model.RoomNotificationRelatedData, err error) {
	var eg errgroup.Group

	var roomId, userId, referRoot int
	var roomTitle, userName string
	eg.Go(func() error {
		row := db.QueryRowx(`
			SELECT
				rooms.id,
				rooms.title,
				characters.id,
				characters.nickname,
				rooms_messages.refer_root
			FROM
				rooms_messages
			JOIN
				rooms ON (rooms_messages.room = rooms.id)
			JOIN
				characters ON (rooms_messages.character = characters.id)
			WHERE
				rooms_messages.id = $1;
		`, messageId)

		return row.Scan(
			&roomId,
			&roomTitle,
			&userId,
			&userName,
			&referRoot,
		)
	})

	repliedWebhooks := make([]string, 0, 16)
	eg.Go(func() error {
		rows, err := db.Queryx(`
			SELECT
				characters.webhook
			FROM
				rooms_messages
			JOIN
				rooms_messages_recipients ON (rooms_messages.id = $1 AND rooms_messages.id = rooms_messages_recipients.message)
			JOIN
				characters ON (rooms_messages_recipients.character = characters.id)
			WHERE
				characters.id      != rooms_messages.character AND
				characters.webhook != ''                       AND
				characters.webhook_replied = true              AND
				characters.deleted_at IS NULL;
		`, messageId)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var webhook string
			err = rows.Scan(&webhook)
			if err != nil {
				return err
			}

			repliedWebhooks = append(repliedWebhooks, webhook)
		}

		return nil
	})

	subscribeWebhooks := make([]string, 0, 16)
	eg.Go(func() error {
		rows, err := db.Queryx(`
			SELECT
				characters.webhook
			FROM
				rooms_messages
			JOIN
				rooms ON (rooms_messages.id = $1 AND rooms.id = rooms_messages.room)
			JOIN
				rooms_message_subscribers ON (rooms.id = rooms_message_subscribers.room)
			JOIN
				characters ON (rooms_message_subscribers.character = characters.id)
			WHERE
				characters.id      != rooms_messages.character AND
				characters.webhook != ''                       AND
				characters.webhook_subscribe = true            AND
				characters.deleted_at IS NULL;
		`, messageId)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var webhook string
			err = rows.Scan(&webhook)
			if err != nil {
				return err
			}

			subscribeWebhooks = append(subscribeWebhooks, webhook)
		}

		return nil
	})

	eg.Go(func() error {
		return db.ExecTx(func(tx *sqlx.Tx) error {
			rows, err := tx.Queryx(`
				INSERT INTO notifications (
					character,
					type
				)

				SELECT
					characters.id,
					'REPLIED'
				FROM
					rooms_messages
				JOIN
					rooms_messages_recipients ON (rooms_messages.id = $1 AND rooms_messages.id = rooms_messages_recipients.message)
				JOIN
					characters ON (rooms_messages_recipients.character = characters.id)
				WHERE
					characters.id != rooms_messages.character AND
					characters.notification_replied = true    AND
					characters.deleted_at IS NULL

				RETURNING
					id;
			`, messageId)
			if err != nil {
				return err
			}
			defer rows.Close()

			type repliedNotificationInserter struct {
				NotificationId int `db:"notification_id"`
				MessageId      int `db:"message_id"`
			}

			inserters := make([]repliedNotificationInserter, 0, 16)
			for rows.Next() {
				var inserter repliedNotificationInserter
				err = rows.Scan(&inserter.NotificationId)
				if err != nil {
					return err
				}
				inserter.MessageId = messageId

				inserters = append(inserters, inserter)
			}

			if 0 < len(inserters) {
				_, err = tx.NamedExec(`
					INSERT INTO notifications_replied_data (
						notification,
						message	
					) VALUES (
						:notification_id,
						:message_id
					)
				`, inserters)
				if err != nil {
					return err
				}
			}

			return nil
		})
	})

	eg.Go(func() error {
		return db.ExecTx(func(tx *sqlx.Tx) error {
			rows, err := tx.Queryx(`
				INSERT INTO notifications (
					character,
					type
				)

				SELECT
					characters.id,
					'SUBSCRIBE'
				FROM
					rooms_messages
				JOIN
					rooms ON (rooms_messages.id = $1 AND rooms.id = rooms_messages.room)
				JOIN
					rooms_message_subscribers ON (rooms.id = rooms_message_subscribers.room)
				JOIN
					characters ON (rooms_message_subscribers.character = characters.id)
				WHERE
					characters.id != rooms_messages.character AND
					characters.notification_subscribe = true  AND
					characters.deleted_at IS NULL

				RETURNING
					id;
			`, messageId)
			if err != nil {
				return err
			}
			defer rows.Close()

			type subscribeNotificationInserter struct {
				NotificationId int `db:"notification_id"`
				MessageId      int `db:"message_id"`
			}

			inserters := make([]subscribeNotificationInserter, 0, 16)
			for rows.Next() {
				var inserter subscribeNotificationInserter
				err = rows.Scan(&inserter.NotificationId)
				if err != nil {
					return err
				}
				inserter.MessageId = messageId

				inserters = append(inserters, inserter)
			}

			if 0 < len(inserters) {
				_, err = tx.NamedExec(`
					INSERT INTO notifications_subscribe_data (
						notification,
						message	
					) VALUES (
						:notification_id,
						:message_id
					)
				`, inserters)
				if err != nil {
					return err
				}
			}

			return nil
		})
	})

	err = eg.Wait()
	if err != nil {
		return nil, err
	}

	return &model.RoomNotificationRelatedData{
		RoomId:            roomId,
		RoomTitle:         roomTitle,
		UserId:            userId,
		UserName:          userName,
		ReferRoot:         referRoot,
		RepliedWebhooks:   repliedWebhooks,
		SubscribeWebhooks: subscribeWebhooks,
	}, nil
}
