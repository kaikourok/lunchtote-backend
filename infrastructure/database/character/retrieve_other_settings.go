package character

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *CharacterRepository) RetrieveOtherSettings(characterId int) (settings *model.CharacterOtherSettingsState, err error) {
	row := db.QueryRowx(`
		SELECT
			characters.webhook,
			characters.webhook_followed,
			characters.webhook_replied,
			characters.webhook_subscribe,
			characters.webhook_mail,
			characters.notification_followed,
			characters.notification_replied,
			characters.notification_subscribe,
			characters.notification_mail,
			characters.email,
			EXISTS (SELECT * FROM characters_twitter WHERE character = $1),
			EXISTS (SELECT * FROM characters_google  WHERE character = $1)
		FROM
			characters
		WHERE
			id = $1;
	`, characterId)

	settings = &model.CharacterOtherSettingsState{}
	err = row.Scan(
		&settings.Webhook.Url,
		&settings.Webhook.Followed,
		&settings.Webhook.Replied,
		&settings.Webhook.Subscribe,
		&settings.Webhook.Mail,
		&settings.Notification.Followed,
		&settings.Notification.Replied,
		&settings.Notification.Subscribe,
		&settings.Notification.Mail,
		&settings.Email,
		&settings.LinkedStates.Twitter,
		&settings.LinkedStates.Google,
	)

	return
}
