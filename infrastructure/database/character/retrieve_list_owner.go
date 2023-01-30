package character

func (db *CharacterRepository) RetrieveListOwner(listId int) (characterId int, err error) {
	row := db.QueryRowx(`
		SELECT
			master
		FROM
			lists
		WHERE
			id = $1;
	`, listId)

	err = row.Scan(&characterId)
	return
}
