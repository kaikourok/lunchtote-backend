package mail

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

func (db *MailRepository) SendMail(userId, targetId int, title, message string) (userName string, targetWebhook string, err error) {
	err = db.ExecTx(func(tx *sqlx.Tx) error {
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
						id = $2 AND administrator = false AND deleted_at IS NULL
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

		row = tx.QueryRowx(`
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
			)
			
			RETURNING
				id;
		`,
			userId,
			targetId,
			title,
			message,
		)

		var mailId int
		err = row.Scan(&mailId)
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
				(SELECT nickname FROM characters WHERE id = $2)
			FROM
				characters
			WHERE
				id = $1;
		`, targetId, userId)

		var isTargetNotificationEnable bool
		err = row.Scan(
			&targetWebhook,
			&isTargetNotificationEnable,
			&userName,
		)
		if err != nil {
			return err
		}

		if isTargetNotificationEnable {
			row = tx.QueryRowx(`
				INSERT INTO notifications (
					character,
					type
				) VALUES (
					$1,
					'MAIL'
				)
				
				RETURNING
					id;
			`, targetId)

			var notificationId int
			err = row.Scan(&notificationId)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`
				INSERT INTO notifications_mail_data (
					notification,
					mail
				) VALUES (
					$1,
					$2
				);
			`, notificationId, mailId)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return
}
