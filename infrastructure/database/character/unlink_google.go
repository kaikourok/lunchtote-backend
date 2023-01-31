package character

func (db *CharacterRepository) UnlinkGoogle(characterId int) error {
	_, err := db.Exec(`
		DELETE FROM
			characters_google
		WHERE
			character = $1;
	`, characterId)

	return err
}
