package diary

func (db *DiaryRepository) ReservePublishDiary(characterId int, title, diary string) error {
	_, err := db.Exec(`
		UPDATE
			characters
		SET
			diary_title = $2,
			diary       = $3
		WHERE
			id = $1;
	`,
		characterId,
		title,
		diary,
	)

	return err
}
