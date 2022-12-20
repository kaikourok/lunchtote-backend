package diary

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *DiaryRepository) RetrieveDiaryPreview(characterId int) (*model.Diary, error) {
	row := db.QueryRowx(`
		SELECT
			characters.id,
			characters.nickname,
			characters.mainicon,
			characters.diary_title,
			characters.diary
		FROM
			characters
		WHERE
			id = $1;
	`, characterId)

	var diary model.Diary
	var diaryTitle, diaryBody *string
	err := row.Scan(
		&diary.Author.Id,
		&diary.Author.Name,
		&diary.Author.Mainicon,
		&diaryTitle,
		&diaryBody,
	)
	if err != nil {
		return nil, err
	}
	if diaryTitle == nil || diaryBody == nil {
		return nil, nil
	}

	diary.Title = *diaryTitle
	diary.Diary = *diaryBody

	return &diary, nil
}
