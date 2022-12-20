package character

func (db *CharacterRepository) CheckUsernameExists(username string) (exists bool, err error) {
	row := db.QueryRowx(`
		SELECT
			EXISTS (
				SELECT
					*
				FROM
					characters
				WHERE
					username = $1
			)
	`, username)

	err = row.Scan(&exists)
	return
}
