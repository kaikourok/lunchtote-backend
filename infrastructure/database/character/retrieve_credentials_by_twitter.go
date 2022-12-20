package character

func (db *CharacterRepository) RetrieveCredentialsByTwitter(twitterId string) (characterId int, notificationToken string, err error) {
	row := db.QueryRowx(`
		SELECT
			characters.id,
			characters.notification_token
		FROM
			characters
		JOIN
			characters_twitter ON (characters.id = characters_twitter.character)
		WHERE
			characters_twitter.twitter_id = $1 AND
			characters.deleted_at IS NULL;
	`, twitterId)

	err = row.Scan(&characterId, &notificationToken)
	if err != nil {
		return 0, "", err
	}

	return characterId, notificationToken, nil
}
