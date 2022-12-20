package character

import (
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) RetrieveCharacterSuggestions(id int, searchText string, excludeOwn bool) (suggestions *model.CharacterSuggestionsData, err error) {
	number, err := strconv.Atoi(searchText)

	var excludeOwnConditionString string
	if excludeOwn {
		excludeOwnConditionString = "id != $1 AND"
	} else {
		excludeOwnConditionString = ""
	}

	var rows *sqlx.Rows
	if err != nil {
		// 数字で検索している場合（Idもしくはキャラクター名検索）
		rows, err = db.Queryx(`
			WITH dismiss_list AS (
				SELECT
					blocker AS dismiss			
				FROM
					blocks
				WHERE
					blocked = $1
				UNION ALL
				SELECT
					blocked AS dismiss
				FROM
					blocks
				WHERE
					blocker = $1
			)
			
			SELECT
				id,
				name
			FROM
				characters
			WHERE
				(
					name     LIKE likequery($2) OR
					nickname LIKE likequery($2)
				) AND
				`+excludeOwnConditionString+`
				deleted_at IS NULL AND
				administrator = false AND
				NOT EXISTS (SELECT * FROM dismiss_list WHERE id = dismiss)
			LIMIT
				20;
		`, id, searchText)
	} else {
		// 数字以外で検索している場合（キャラクター名検索）
		rows, err = db.Queryx(`
			WITH dismiss_list AS (
				SELECT
					blocker AS dismiss			
				FROM
					blocks
				WHERE
					blocked = $1
				UNION ALL
				SELECT
					blocked AS dismiss
				FROM
					blocks
				WHERE
					blocker = $1
			)
			
			SELECT
				id,
				name
			FROM
				characters
			WHERE
				(
					id = $2 OR
					name     LIKE likequery($3) OR
					nickname LIKE likequery($3)
				) AND
				`+excludeOwnConditionString+`
				deleted_at IS NULL AND
				administrator = false AND
				NOT EXISTS (SELECT * FROM dismiss_list WHERE id = dismiss)
			LIMIT
				20;
		`, id, number, searchText)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	suggestionsData := make(model.CharacterSuggestionsData, 0, 20)
	for rows.Next() {
		var character model.CharacterSuggestionData
		err = rows.Scan(
			&character.Id,
			&character.Name,
		)
		if err != nil {
			return nil, err
		}

		suggestionsData = append(suggestionsData, character)
	}

	return &suggestionsData, nil
}
