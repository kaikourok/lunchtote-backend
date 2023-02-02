package mail

import (
	"github.com/jmoiron/sqlx"
)

func (db *MailRepository) SendAdministratorMail(targetId *int, name, title, message string) (webhooks []string, err error) {
	webhooks = make([]string, 0, 2048)

	err = db.ExecTx(func(tx *sqlx.Tx) error {
		if targetId == nil {
			row := tx.QueryRowx(`
				INSERT INTO mails (
					name,
					title,
					message,
					receiver
				) VALUES (
					$1,
					$2,
					$3,
					$4
				)
				
				RETURNING
					id;
			`,
				name,
				title,
				message,
				targetId,
			)

			var mailId int
			err := row.Scan(&mailId)
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
					deleted_at IS NOT NULL
				FROM
					characters
				WHERE
					id = $1;
			`, targetId)

			var webhook string
			var isTargetNotificationEnable, isTargetDeleted bool
			err = row.Scan(
				webhook,
				isTargetNotificationEnable,
				isTargetDeleted,
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

			if !isTargetDeleted && webhook != "" {
				webhooks = append(webhooks, webhook)
			}
		} else {
			row := tx.QueryRowx(`
				INSERT INTO mass_mails (
					name,
					title,
					message
				) VALUES (
					$1,
					$2,
					$3
				)
				
				RETURNING
					id;
			`,
				name,
				title,
				message,
			)

			var massMailId int
			err = row.Scan(&massMailId)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`
				INSERT INTO mails (
					name,
					title,
					message,
					receiver
				)
				SELECT
					$1,
					$2,
					$3,
					id
				FROM
					characters;
			`,
				name,
				title,
				message,
			)
			if err != nil {
				return err
			}

			rows, err := tx.Queryx(`
				SELECT
					webhook
				FROM
					characters
				WHERE
					webhook      != ''   AND
					webhook_mail  = true AND
					deleted_at IS NULL;
			`)
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

				webhooks = append(webhooks, webhook)
			}

			_, err = tx.Exec(`
				INSERT INTO notifications (
					character,
					message
				)
				SELECT
					id,
					$1
				FROM
					characters
				WHERE
					notification_mail = true;
			`)
			if err != nil {
				return err
			}

			rows, err = tx.Queryx(`
				INSERT INTO (
					character,
					type
				)

				SELECT
					characters.id,
					'MASS_MAIL'
				FROM
					characters
				WHERE
					characters.notification_mail = true AND
					characters.deleted_at IS NULL
					
				RETURNING
					id;
			`)
			if err != nil {
				return err
			}
			defer rows.Close()

			type notificationIdInsertStruct struct {
				NotificationId int `db:"id"`
				MassMailId     int `db:"mail_id"`
			}

			notificationIds := make([]notificationIdInsertStruct, 0, 2048)
			for rows.Next() {
				var notificationId notificationIdInsertStruct
				err = rows.Scan(&notificationId.NotificationId)
				if err != nil {
					return err
				}
				notificationId.MassMailId = massMailId

				notificationIds = append(notificationIds, notificationId)
			}

			_, err = db.NamedExec(`
				INSERT INTO notifications_mass_mail_data (
					notification,
					mail
				) VALUES (
					:id,
					:mail_id
				)
			`, notificationIds)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return
}
