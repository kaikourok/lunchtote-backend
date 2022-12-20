package diary

func (db *DiaryRepository) ClearReservedDiary(characterId int) error {
	_, err := db.Exec(`
		UPDATE
			characters
		SET
			diary_title = NULL,
			diary       = NULL
		WHERE
			id = $1;
	`, characterId)
	if err != nil {
		return err
	}

	return nil
}
