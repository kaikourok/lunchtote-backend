package character

import (
	"errors"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) RetrieveListSuggestions(characterId int, searchText string, listId int) (suggestions *model.CharacterSuggestionsData, err error) {
	row := db.QueryRowx(`
		SELECT
			master
		FROM
			lists
		WHERE
			id = $1;
	`, listId)

	var listMaster int
	err = row.Scan(&listMaster)
	if err != nil {
		return nil, err
	}

	if characterId != listMaster {
		err := errors.New("リストの管理者ではありません")
		return nil, err
	}

	number, err := strconv.Atoi(searchText)

	var rows *sqlx.Rows
	if err != nil {
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
				NOT EXISTS (SELECT * FROM lists_characters WHERE list = $3 AND character = characters.id) AND
				deleted_at IS NULL AND
				administrator = false AND
				NOT EXISTS (SELECT * FROM dismiss_list WHERE id = dismiss)
			LIMIT
				20;
		`, characterId, searchText, listId)
	} else {
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
				NOT EXISTS (SELECT * FROM lists_characters WHERE list = $4 AND character = characters.id) AND
				deleted_at IS NULL AND
				administrator = false AND
				NOT EXISTS (SELECT * FROM dismiss_list WHERE id = dismiss)
			LIMIT
				20;
		`, characterId, number, searchText, listId)
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
