package character

import (
	"encoding/json"
	"errors"

	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func unmarshalValueJson[T any](valueJson []byte) (value T, err error) {
	err = json.Unmarshal(valueJson, &value)
	return
}

func (db *CharacterRepository) RetrieveNotifications(id, start, number int) (notifications []model.Notification, isContinue bool, err error) {
	rows, err := db.Queryx(`
		WITH target_notifications AS (
			SELECT
				id,
				type,
				notificated_at AS timestamp
			FROM
				notifications
			WHERE
				character = $1 AND id > $2
			ORDER BY
				notificated_at DESC
			LIMIT
				$3
		)

		SELECT
			valued_notifications.id,
			valued_notifications.type,
			valued_notifications.timestamp,
			valued_notifications.value
		FROM
			(
				SELECT
					target_notifications.id,
					target_notifications.type,
					target_notifications.timestamp,
					JSON_BUILD_OBJECT(
						'character', JSON_BUILD_OBJECT(
							'id',       characters.id,
							'name',     characters.nickname,
							'mainicon', characters.mainicon
						)
					) AS value
				FROM
					target_notifications
				JOIN
					notifications_followed_data ON (target_notifications.type = 'FOLLOWED' AND target_notifications.id = notifications_followed_data.notification)
				JOIN
					characters ON (notifications_followed_data.followed_by = characters.id)

				UNION ALL

				SELECT
					target_notifications.id,
					target_notifications.type,
					target_notifications.timestamp,
					JSON_BUILD_OBJECT(
						'character', JSON_BUILD_OBJECT(
							'id',       characters.id,
							'name',     characters.nickname,
							'mainicon', characters.mainicon
						),
						'room', JSON_BUILD_OBJECT(
							'id',    rooms.id,
							'title', rooms.title
						),
						'message', JSON_BUILD_OBJECT(
							'referRoot', rooms_messages.refer_root,
							'message',   rooms_messages.message
						)
					) AS value
				FROM
					target_notifications
				JOIN
					notifications_replied_data ON (target_notifications.type = 'REPLIED' AND target_notifications.id = notifications_replied_data.notification)
				JOIN
					rooms_messages ON (notifications_replied_data.message = rooms_messages.id)
				JOIN
					rooms ON (rooms_messages.room = rooms.id)
				JOIN
					characters ON (rooms_messages.character = characters.id)

				UNION ALL
				
				SELECT
					target_notifications.id,
					target_notifications.type,
					target_notifications.timestamp,
					JSON_BUILD_OBJECT(
						'character', JSON_BUILD_OBJECT(
							'id',       characters.id,
							'name',     characters.nickname,
							'mainicon', characters.mainicon
						),
						'room', JSON_BUILD_OBJECT(
							'id',    rooms.id,
							'title', rooms.title
						),
						'message', JSON_BUILD_OBJECT(
							'message',   rooms_messages.message
						)
					) AS value
				FROM
					target_notifications
				JOIN
					notifications_subscribe_data ON (target_notifications.type = 'SUBSCRIBE' AND target_notifications.id = notifications_subscribe_data.notification)
				JOIN
					rooms_messages ON (notifications_subscribe_data.message = rooms_messages.id)
				JOIN
					rooms ON (rooms_messages.room = rooms.id)
				JOIN
					characters ON (rooms_messages.character = characters.id)

				UNION ALL
				
				SELECT
					target_notifications.id,
					target_notifications.type,
					target_notifications.timestamp,
					JSON_BUILD_OBJECT(
						'character', JSON_BUILD_OBJECT(
							'id',       characters.id,
							'name',     characters.nickname,
							'mainicon', characters.mainicon
						),
						'room', JSON_BUILD_OBJECT(
							'id',    rooms.id,
							'title', rooms.title
						)
					) AS value
				FROM
					target_notifications
				JOIN
					notifications_new_member_data ON (target_notifications.type = 'NEW_MEMBER' AND target_notifications.id = notifications_new_member_data.notification)
				JOIN
					rooms ON (notifications_new_member_data.room = rooms.id)
				JOIN
					characters ON (notifications_new_member_data.character = characters.id)

				UNION ALL
				
				SELECT
					target_notifications.id,
					target_notifications.type,
					target_notifications.timestamp,
					JSON_BUILD_OBJECT(
						'character', JSON_BUILD_OBJECT(
							'id',       characters.id,
							'name',     characters.nickname,
							'mainicon', characters.mainicon
						),
						'mail', JSON_BUILD_OBJECT(
							'id',    mails.id,
							'title', mails.title
						)
					) AS value
				FROM
					target_notifications
				JOIN
					notifications_mail_data ON (target_notifications.type = 'MAIL' AND target_notifications.id = notifications_mail_data.notification)
				JOIN
					mails ON (notifications_mail_data.mail = mails.id)
				JOIN
					characters ON (mails.sender = characters.id)

				UNION ALL
				
				SELECT
					target_notifications.id,
					target_notifications.type,
					target_notifications.timestamp,
					JSON_BUILD_OBJECT(
						'mail', JSON_BUILD_OBJECT(
							'id',    mass_mails.id,
							'title', mass_mails.title
						)
					) AS value
				FROM
					target_notifications
				JOIN
					notifications_mass_mail_data ON (target_notifications.type = 'MASS_MAIL' AND target_notifications.id = notifications_mass_mail_data.notification)
				JOIN
					mass_mails ON (notifications_mass_mail_data.mail = mass_mails.id)
			) valued_notifications
		ORDER BY
			valued_notifications.timestamp DESC;
	`, id, start, number+1)
	if err != nil {
		return
	}
	defer rows.Close()

	notifications = make([]model.Notification, 0, number+1)
	for rows.Next() {
		var base model.NotificationBase
		var valueJson string
		err = rows.Scan(
			&base.Id,
			&base.Type,
			&base.Timestamp,
			&valueJson,
		)
		if err != nil {
			return
		}

		var notification model.Notification
		switch base.Type {
		case model.NotificationTypeFollowed:
			value, err := unmarshalValueJson[model.FollowedNotificationValue]([]byte(valueJson))
			if err != nil {
				return nil, false, err
			}
			notification = model.FollowedNotification{
				NotificationBase:          base,
				FollowedNotificationValue: value,
			}
		case model.NotificationTypeReplied:
			value, err := unmarshalValueJson[model.RepliedNotificationValue]([]byte(valueJson))
			if err != nil {
				return nil, false, err
			}
			notification = model.RepliedNotification{
				NotificationBase:         base,
				RepliedNotificationValue: value,
			}
		case model.NotificationTypeSubscribe:
			value, err := unmarshalValueJson[model.SubscribeNotificationValue]([]byte(valueJson))
			if err != nil {
				return nil, false, err
			}
			notification = model.SubscribeNotification{
				NotificationBase:           base,
				SubscribeNotificationValue: value,
			}
		case model.NotificationTypeNewMember:
			value, err := unmarshalValueJson[model.NewMemberNotificationValue]([]byte(valueJson))
			if err != nil {
				return nil, false, err
			}
			notification = model.NewMemberNotification{
				NotificationBase:           base,
				NewMemberNotificationValue: value,
			}
		case model.NotificationTypeMail:
			value, err := unmarshalValueJson[model.MailNotificationValue]([]byte(valueJson))
			if err != nil {
				return nil, false, err
			}
			notification = model.MailNotification{
				NotificationBase:      base,
				MailNotificationValue: value,
			}
		case model.NotificationTypeMassMail:
			value, err := unmarshalValueJson[model.MassMailNotificationValue]([]byte(valueJson))
			if err != nil {
				return nil, false, err
			}
			notification = model.MassMailNotification{
				NotificationBase:          base,
				MassMailNotificationValue: value,
			}
		default:
			return nil, false, errors.New("不明な通知タイプです:" + base.Type)
		}

		notifications = append(notifications, notification)
	}

	if len(notifications) == number+1 {
		notifications = notifications[:len(notifications)-1]
		isContinue = true
	}

	return
}
