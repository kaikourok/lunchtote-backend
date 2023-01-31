package character

func (db *CharacterRepository) UnlinkTwitter(characterId int) error {
	_, err := db.Exec(`
		DELETE FROM
			characters_twitter
		WHERE
			character = $1;
	`, characterId)

	return err
}
