package character

func (db *CharacterRepository) RetrieveCharacterNickname(characterId int) (nickname string, err error) {
	row := db.QueryRowx(`
		SELECT
			nickname
		FROM
			characters
		WHERE
			id = $1;
	`, characterId)

	err = row.Scan(&nickname)
	return
}
