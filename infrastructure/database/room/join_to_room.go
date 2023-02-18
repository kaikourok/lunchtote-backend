package room

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *RoomRepository) JoinToRoom(targetId, roomId int) (room *model.RoomOverview, targetName string, newMemberWebhooks []string, err error) {
	row := db.QueryRowx(`
		SELECT
			EXISTS (
				SELECT
					*
				FROM
					rooms_members
				WHERE
					room   = $1 AND
					member = $2 
			)
		;
	`, roomId, targetId)

	var isAlreadyExists bool
	err = row.Scan(&isAlreadyExists)
	if err != nil {
		return nil, "", nil, err
	}

	if isAlreadyExists {
		return nil, "", nil, errors.New("すでにメンバーになっているキャラクターです")
	}

	err = db.ExecTx(func(tx *sqlx.Tx) error {
		rows, err := tx.Queryx(`
			WITH character AS (
				SELECT
					nickname
				FROM
					characters
				WHERE
					id = $1 AND deleted_at IS NULL
			), room AS (
				SELECT
					id,
					title
				FROM
					rooms
				WHERE
					id = $2 AND deleted_at IS NULL
			)

			SELECT
				(SELECT nickname FROM character),
				(SELECT id       FROM room),
				(SELECT title    FROM room),
				characters.webhook
			FROM
				rooms_new_member_subscribers
			JOIN
				characters ON (rooms_new_member_subscribers.room = $2 AND rooms_new_member_subscribers.character = characters.id)
			WHERE
				characters.webhook != ''             AND
				characters.webhook_new_member = true AND
				characters.deleted_at IS NULL;
		`, targetId, roomId)
		if err != nil {
			return err
		}
		defer rows.Close()

		room = &model.RoomOverview{}
		newMemberWebhooks = make([]string, 0, 16)
		for rows.Next() {
			var webhook string
			err = rows.Scan(
				&targetName,
				&room.Id,
				&room.Title,
				&webhook,
			)
			if err != nil {
				return err
			}

			newMemberWebhooks = append(newMemberWebhooks, webhook)
		}

		rows, err = tx.Queryx(`
			INSERT INTO notifications (
				character,
				type
			)
			
			SELECT
				characters.id,
				'NEW_MEMBER'
			FROM
				rooms_new_member_subscribers
			JOIN
				characters ON (rooms_new_member_subscribers.room = $1 AND rooms_new_member_subscribers.character = characters.id)
			WHERE
				characters.notification_new_member = true AND
				characters.deleted_at IS NULL
				
			RETURNING
				id;
		`, roomId)
		if err != nil {
			return err
		}
		defer rows.Close()

		type newMemberNotificationInserter struct {
			NotificationId int `db:"notification_id"`
			CharacterId    int `db:"character_id"`
			RoomId         int `db:"room_id"`
		}

		inserters := make([]newMemberNotificationInserter, 0, 16)
		for rows.Next() {
			var inserter newMemberNotificationInserter
			err = rows.Scan(&inserter.NotificationId)
			if err != nil {
				return err
			}
			inserter.CharacterId = targetId
			inserter.RoomId = roomId

			inserters = append(inserters, inserter)
		}

		if 0 < len(inserters) {
			_, err = tx.NamedExec(`
				INSERT INTO notifications_new_member_data (
					notification,
					room,
					character
				) VALUES (
					:notification_id,
					:room_id,
					:character_id
				)
			`, inserters)
			if err != nil {
				return err
			}
		}

		_, err = tx.Exec(`
			INSERT INTO rooms_members (
				room,
				member,
				write,
				ban,
				invite,
				use_reply,
				use_secret,
				delete_other_message,
				create_children_room,
				type
			)

			SELECT
				$1,
				$2,
				write,
				ban,
				invite,
				use_reply,
				use_secret,
				delete_other_message,
				create_children_room,
				'MEMBER'
			FROM
				rooms_roles
			WHERE
				rooms_roles.room = $1 AND
				rooms_roles.type = 'DEFAULT';
		`, roomId, targetId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO rooms_roles_members (
				role,
				character
			)

			SELECT
				id,
				$2
			FROM
				rooms_roles
			WHERE
				rooms_roles.room = $1 AND
				rooms_roles.type = 'DEFAULT';
		`, roomId, targetId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE
				rooms
			SET
				members_count = members_count + 1
			WHERE
				id = $1;
		`, roomId)
		if err != nil {
			return err
		}

		return nil
	})

	return
}
