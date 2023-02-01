package character

func (db *CharacterRepository) UpdateNotificationChecked(characterId int) error {
	_, err := db.Exec(`
		UPDATE
			characters
		SET
			notification_last_checked_at = CURRENT_TIMESTAMP
		WHERE
			id = $1;
	`, characterId)

	return err
}
