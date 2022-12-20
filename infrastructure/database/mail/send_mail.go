package mail

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

func (db *MailRepository) SendMail(userId, targetId int, title, message string) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		row := tx.QueryRowx(`
			SELECT
				EXISTS (
					SELECT
						*
					FROM
						blocks
					WHERE
						(blocker = $1 AND blocked = $2) OR
						(blocked = $1 AND blocker = $2)
				),
				EXISTS (
					SELECT
						*
					FROM
						characters
					WHERE
						id = $1 AND administrator = false AND deleted_at IS NULL
				);
		`, userId, targetId)

		var isForbidden, isFound bool
		err := row.Scan(&isForbidden, &isFound)

		if err != nil {
			return err
		}

		if isForbidden {
			return errors.New("ブロックしている、もしくはブロックされています")
		}

		if !isFound {
			return errors.New("対象のキャラクターが存在しません")
		}

		_, err = tx.Exec(`
			INSERT INTO mails (
				sender,
				receiver,
				title,
				message
			) VALUES (
				$1,
				$2,
				$3,
				$4
			);
		`,
			userId,
			targetId,
			title,
			message,
		)
		if err != nil {
			return err
		}

		row = tx.QueryRowx(`
			SELECT
				CASE webhook_mail
					WHEN true  THEN webhook
					WHEN false THEN ''
				END,
				notification_mail,
				deleted_at IS NOT NULL,
				(SELECT nickname FROM characters WHERE id = $2)
			FROM
				characters
			WHERE
				id = $1;
		`, targetId, userId)

		var webhook, userName string
		var isTargetNotificationEnable, isTargetDeleted bool
		err = row.Scan(
			&webhook,
			&isTargetNotificationEnable,
			&isTargetDeleted,
			&userName,
		)
		if err != nil {
			return err
		}

		return nil
	})

	/* TODO
		notificationMessage := strings.ReplaceAll(
			strings.ReplaceAll(config.GetString("notification.mail-template"), "{entry-number}", utils.ConvertCharacterIdToText(userId)),
			"{name}", userName,
		)

		if isTargetNotificationEnable {
			_, err = tx.Exec(`
				INSERT INTO notifications (
					type,
					character,
					message,
					icon,
					value,
					detail
				) VALUES (
					'FOLLOWED',
					$1,
					$2,
					(SELECT mainicon FROM characters WHERE id = $3),
					cast($3 as TEXT),
					$4
				);
			`, payload.Target, notificationMessage, userId, payload.Message)
			if err != nil {
				log.Println(err)
				return err
			}
		}

		if !isTargetDeleted {
			if webhook != "" {
				go notification.SendWebhook(webhook, notificationMessage)
			}
			if isTargetNotificationEnable {
				go notification.SendNotification([]int{payload.Target}, notificationMessage)
			}
		}

		return nil
	})
	*/
}
