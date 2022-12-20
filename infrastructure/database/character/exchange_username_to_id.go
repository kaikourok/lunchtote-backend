package character

func (db *CharacterRepository) ExchangeUsernameToId(username string) (id int, err error) {
	row := db.QueryRowx(`
		SELECT
			id
		FROM
			characters
		WHERE
			username = $1
	`, username)

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
