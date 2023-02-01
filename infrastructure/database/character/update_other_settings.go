package character

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *CharacterRepository) UpdateOtherSettings(characterId int, settings *model.CharacterOtherSettings) error {
	_, err := db.Exec(`
		UPDATE
			characters
		SET
			webhook                 = $2,
			webhook_followed        = $3,
			webhook_replied         = $4,
			webhook_subscribe       = $5,
			webhook_new_member      = $6,
			webhook_mail            = $7,
			notification_followed   = $8,
			notification_replied    = $9,
			notification_subscribe  = $10,
			notification_new_member = $11,
			notification_mail       = $12
		WHERE
			id = $1;
	`,
		characterId,
		&settings.Webhook.Url,
		&settings.Webhook.Followed,
		&settings.Webhook.Replied,
		&settings.Webhook.Subscribe,
		&settings.Webhook.NewMember,
		&settings.Webhook.Mail,
		&settings.Notification.Followed,
		&settings.Notification.Replied,
		&settings.Notification.Subscribe,
		&settings.Notification.NewMember,
		&settings.Notification.Mail,
	)

	return err
}
