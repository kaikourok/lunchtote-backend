package character

func (db *CharacterRepository) CreateLayeringGroup(characterId int, name string) (id int, err error) {
	row := db.QueryRowx(`
		INSERT INTO characters_icon_layering_groups (
			character,
			name
		) VALUES (
			$1,
			$2
		)
		RETURNING
			id;
	`, characterId, name)

	err = row.Scan(&id)
	return
}
