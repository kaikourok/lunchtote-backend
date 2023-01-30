package character

func (db *CharacterRepository) RenameList(listId int, newName string) error {
	_, err := db.Exec(`
		UPDATE
			lists
		SET
			name = $2		
		WHERE
			id = $1;
	`, listId, newName)

	return err
}
