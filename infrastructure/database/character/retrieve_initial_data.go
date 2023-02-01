package character

func (db *CharacterRepository) RetrieveInitialData(id int) (existsUnreadNotification, existsUnreadMail bool, err error) {
	row := db.QueryRowx(`
		SELECT
			EXISTS (
				SELECT
					*
				FROM
					notifications
				WHERE
					character = characters.id AND
					notificated_at > (SELECT notification_last_checked_at FROM characters WHERE id = $1)
				),
			EXISTS (
				SELECT
					*
				FROM
					mails
				WHERE
					receiver = characters.id AND
					read     = false AND
					deleted_at IS NULL
				)
		FROM
			characters
		WHERE
			characters.id = $1 AND
			characters.deleted_at IS NULL;
	`, id)

	err = row.Scan(
		&existsUnreadNotification,
		&existsUnreadMail,
	)

	return
}
