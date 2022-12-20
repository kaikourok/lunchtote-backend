package character

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *CharacterRepository) RetrieveNotifications(id, start, number int) (notifications *[]model.Notification, isContinue bool, err error) {
	rows, err := db.Queryx(`
		SELECT
			type,
			icon,
			message,
			detail,
			value,
			notificated_at
		FROM
			notifications
		WHERE
			character = $1 AND id > $2
		ORDER BY
			notificated_at DESC
		LIMIT
			$3;
	`, id, start, number+1)
	if err != nil {
		return
	}
	defer rows.Close()

	notificationsSlice := make([]model.Notification, 0, number+1)
	for rows.Next() {
		var notification model.Notification
		err = rows.Scan(
			&notification.Type,
			&notification.Icon,
			&notification.Message,
			&notification.Detail,
			&notification.Value,
			&notification.Timestamp,
		)
		if err != nil {
			return
		}

		notificationsSlice = append(notificationsSlice, notification)
	}

	if len(notificationsSlice) == number+1 {
		notificationsSlice = notificationsSlice[:len(notificationsSlice)-1]
		isContinue = true
	}

	return &notificationsSlice, isContinue, nil
}
