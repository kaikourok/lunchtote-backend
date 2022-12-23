package diary

import (
	"errors"

	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *DiaryRepository) RetrieveDiary(characterId *int, targetId, nth int) (*model.Diary, error) {
	if characterId != nil {
		row := db.QueryRowx(`
			SELECT
				EXISTS (
					SELECT
						*
					FROM
						blocks
					WHERE
						(blocker = $1 AND blocked = $2) OR
						(blocker = $2 AND blocked = $1)
				);
		`)

		var isInvalid bool
		err := row.Scan(&isInvalid)
		if err != nil {
			return nil, err
		}

		if isInvalid {
			return nil, errors.New("ブロックしている、もしくはブロックされているため取得できません")
		}
	}

	row := db.QueryRowx(`
		SELECT
			characters.id,
			characters.nickname,
			characters.mainicon,
			diaries.title,
			diaries.diary,
			diaries.nth
		FROM
			diaries
		JOIN
			characters ON (diaries.character = characters.character_id)
		WHERE
			diaries.character = $1 AND diaries.nth = $2;
	`, targetId, nth)

	var diary model.Diary
	err := row.Scan(
		&diary.Author.Id,
		&diary.Author.Name,
		&diary.Author.Mainicon,
		&diary.Title,
		&diary.Diary,
		&diary.Nth,
	)
	if err != nil {
		return nil, err
	}

	rows, err := db.Queryx(`
		SELECT
			diaries.nth
		FROM
			diaries
		WHERE
			diaries.character = $1
		ORDER BY
			diaries.nth;
	`, targetId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	diaries := make([]int, 0, 64)

	for rows.Next() {
		var diaryIndex int
		err = rows.Scan(&diaryIndex)
		if err != nil {
			return nil, err
		}

		diaries = append(diaries, diaryIndex)
	}

	diary.ExistingDiaries = diaries

	return &diary, nil
}
