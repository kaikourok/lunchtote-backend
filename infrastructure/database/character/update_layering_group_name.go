package character

func (db *CharacterRepository) UpdateLayeringGroupName(characterId, groupId int, name string) error {
	_, err := db.Exec(`
		UPDATE
			characters_icon_layering_groups
		SET
			name = $3
		WHERE
			character = $1 AND id = $2;
	`, characterId, groupId, name)

	return err
}
