package character

func (db *CharacterRepository) ExchangeNotificationTokenToId(token string) (id int, err error) {
	row := db.QueryRowx(`
		SELECT
			id
		FROM
			characters
		WHERE
			notification_token = $1
	`, token)

	err = row.Scan(&id)
	return
}
