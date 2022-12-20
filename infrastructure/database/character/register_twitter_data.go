package character

func (db *CharacterRepository) RegisterTwitterData(characterId int, twitterId string) error {
	_, err := db.Exec(`
		INSERT INTO characters_twitter (
      character,
      twitter_id,
		) VALUES (
			$1,
			$2
		)
		ON CONFLICT (character) DO UPDATE SET
      twitter_id = excluded.twitter_id;
	`,
		characterId,
		twitterId,
	)
	if err != nil {
		return err
	}

	return nil
}
