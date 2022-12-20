package character

func (db *CharacterRepository) ExchangeEmailToId(email string) (id int, err error) {
	row := db.QueryRowx(`
		SELECT
			id
		FROM
			characters
		WHERE
			email = $1
	`, email)

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
