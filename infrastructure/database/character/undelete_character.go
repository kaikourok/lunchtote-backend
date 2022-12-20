package character

func (db *CharacterRepository) UndeleteCharacter(id int) error {
	_, err := db.Exec(`
		UPDATE
			characters
		SET
			deleted_at = NULL
		WHERE
			id = $1;
	`)

	return err
}
