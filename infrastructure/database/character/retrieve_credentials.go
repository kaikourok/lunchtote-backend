package character

func (db *CharacterRepository) RetrieveCredentials(id int) (password, notificationToken string, isAdministrator bool, err error) {
	row := db.QueryRowx(`
		SELECT
			password,
			notification_token,
			administrator
		FROM
			characters
		WHERE
			id = $1 AND
			deleted_at IS NULL;
	`, id)

	err = row.Scan(&password, &notificationToken, &isAdministrator)
	if err != nil {
		return "", "", false, err
	}

	return
}
