package diary

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/lib/pq"
)

func (db *DiaryRepository) RetrieveDiaryEditData(characterId int) (*model.DiaryEditData, error) {

	row := db.QueryRowx(`
		SELECT
			COALESCE(diary_title, ''),
			COALESCE(diary,       ''),
			ARRAY(SELECT path FROM characters_icons WHERE character = $1)
		FROM
			characters
		WHERE
			id = $1;
	`, characterId)

	var diary model.DiaryEditData
	err := row.Scan(
		&diary.Title,
		&diary.Diary,
		pq.Array(&diary.SelectableIcons),
	)
	if err != nil {
		return nil, err
	}

	return &diary, nil
}
