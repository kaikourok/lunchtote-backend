package character

func (db *CharacterRepository) UnregisterEmail(characterId int) error {
	_, err := db.Exec(`
		UPDATE
			characters
		SET
			email = NULL
		WHERE
			id = $1;
	`, characterId)

	return err
}
