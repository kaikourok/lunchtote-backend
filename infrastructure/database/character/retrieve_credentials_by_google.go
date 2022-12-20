package character

func (db *CharacterRepository) RetrieveCredentialsByGoogle(googleId string) (characterId int, notificationToken string, err error) {
	row := db.QueryRowx(`
		SELECT
			characters.id,
			characters.notification_token
		FROM
			characters
		JOIN
			characters_google ON (characters.id = characters_google.character)
		WHERE
			characters_google.google_id = $1 AND
			characters.deleted_at IS NULL;
	`, googleId)

	err = row.Scan(&characterId, &notificationToken)
	if err != nil {
		return 0, "", err
	}

	return characterId, notificationToken, nil
}
