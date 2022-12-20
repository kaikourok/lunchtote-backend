package mail

import (
	"github.com/jmoiron/sqlx"
)

func (db *MailRepository) SendAdministratorMail(targetId *int, name, title, message string) error {
	//notificationMessage := strings.ReplaceAll(config.GetString("notification.administrator-mail-template"), "{name}", name)

	webhooks := make([]string, 0, 2048)
	notificationTargets := make([]int, 0, 2048)

	return db.ExecTx(func(tx *sqlx.Tx) error {
		if targetId == nil {
			_, err := tx.Exec(`
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
				);
			`,
				name,
				title,
				message,
				targetId,
			)
			if err != nil {
				return err
			}

			row := tx.QueryRowx(`
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
				_, err = tx.Exec(`
					INSERT INTO notifications (
						character,
						message
					) VALUES (
						$1,
						$2
					);
				`, targetId, message)
				if err != nil {
					return err
				}
			}

			if !isTargetDeleted {
				if webhook != "" {
					webhooks = append(webhooks, webhook)
				}
				if isTargetNotificationEnable {
					notificationTargets = append(notificationTargets, *targetId)
				}
			}
		} else {
			_, err := tx.Exec(`
				INSERT INTO mails_everyone (
					name,
					title,
					message
				) VALUES (
					$1,
					$2,
					$3
				);
			`,
				name,
				title,
				message,
			)
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
					webhook      = ''   AND
					webhook_mail = true AND
					deleted_at IS NULL;
			`)
			if err != nil {
				return err
			}

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
				SELECT
					id
				FROM
					characters
				WHERE
					notification_mail = true AND deleted_at IS NULL;
			`)
			if err != nil {
				return err
			}

			for rows.Next() {
				var characterId int
				err = rows.Scan(&characterId)
				if err != nil {
					return err
				}

				notificationTargets = append(notificationTargets, characterId)
			}
		}

		return nil
	})

	/* TODO
	go (func() {
		for _, webhook := range webhooks {
			notification.SendWebhook(webhook, notificationMessage)
		}
	})()

	go (func() {
		for _, target := range notificationTargets {
			notification.SendNotification([]int{target}, notificationMessage)
		}
	})()
	*/
}
