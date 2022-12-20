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
					character = $1 AND
					notificated_at > (SELECT notification_last_checked_at FROM characters WHERE id = $1)
				),
			EXISTS (
				SELECT
					*
				FROM
					mails
				WHERE
					receiver = $1    AND
					read     = false AND
					deleted_at IS NULL
				);
	`, id)

	err = row.Scan(
		&existsUnreadNotification,
		&existsUnreadMail,
	)

	return
}
