package character

func (db *CharacterRepository) DeleteCharacter(id int) error {
	_, err := db.Exec(`
		UPDATE
			characters
		SET
			deleted_at = CURRENT_TIMESTAMP
		WHERE
			id = $1;
	`)

	return err
}
