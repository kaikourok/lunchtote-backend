package character

func (db *CharacterRepository) CreateList(characterId int, name string) (listId int, err error) {
	row := db.QueryRowx(`
		INSERT INTO lists (
			master,
			name
		) VALUES (
			$1,
			$2
		)
		RETURNING
			id;
	`, characterId, name)

	err = row.Scan(&listId)
	return
}
