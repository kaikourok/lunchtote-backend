package character

func (db *CharacterRepository) RegisterGoogleData(characterId int, googleId string) error {
	_, err := db.Exec(`
		INSERT INTO characters_google (
      character,
      google_id
		) VALUES (
			$1,
			$2
		)
		ON CONFLICT (character) DO UPDATE SET
      google_id = excluded.google_id;
	`,
		characterId,
		googleId,
	)
	if err != nil {
		return err
	}

	return nil
}
