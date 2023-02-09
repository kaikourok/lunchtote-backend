package diary

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *DiaryRepository) RetrieveDiaryOverviews(characterId, nth int) (diaries *[]model.DiaryOverview, err error) {
	rows, err := db.Queryx(`
		SELECT
			characters.id,
			characters.nickname,
			characters.mainicon,
			diaries.title
		FROM
			diaries
		JOIN
			characters ON (diaries.character = characters.id)
		WHERE
			diaries.nth = $2 AND
			diaries.character = $1 OR EXISTS (SELECT * FROM follows WHERE follower = $1 AND followed = diaries.character)
		ORDER BY
			characters.id;
	`, characterId, nth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	diaryOverviews := make([]model.DiaryOverview, 0, 128)

	for rows.Next() {
		var diary model.DiaryOverview
		err = rows.Scan(
			&diary.Author.Id,
			&diary.Author.Name,
			&diary.Author.Mainicon,
			&diary.Title,
		)
		if err != nil {
			return nil, err
		}

		diaryOverviews = append(diaryOverviews, diary)
	}

	return &diaryOverviews, nil
}
