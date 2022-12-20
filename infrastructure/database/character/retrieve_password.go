package character

func (db *CharacterRepository) RetrievePassword(id int) (password string, err error) {
	row := db.QueryRowx(`
		SELECT
			password
		FROM
			characters
		WHERE
			id = $1 AND deleted_at IS NULL;
	`, id)

	err = row.Scan(&password)
	return
}
