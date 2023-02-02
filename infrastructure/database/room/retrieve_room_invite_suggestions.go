package room

import (
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *RoomRepository) RetrieveRoomInviteSuggestions(characterId int, searchText string, roomId int) (suggestions *model.CharacterSuggestionsData, err error) {
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
				deleted_at IS NULL    AND
				administrator = false AND
				NOT EXISTS (SELECT * FROM rooms_members            WHERE room = $3 AND member  = characters.id) AND
				NOT EXISTS (SELECT * FROM rooms_banned_characters  WHERE room = $3 AND banned  = characters.id) AND
				NOT EXISTS (SELECT * FROM rooms_invited_characters WHERE room = $3 AND invited = characters.id) AND
				NOT EXISTS (SELECT * FROM dismiss_list WHERE id = dismiss)
			LIMIT
				20;
		`, characterId, searchText, roomId)
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
				deleted_at IS NULL AND
				administrator = false AND
				NOT EXISTS (SELECT * FROM rooms_members            WHERE room = $4 AND member  = characters.id) AND
				NOT EXISTS (SELECT * FROM rooms_banned_characters  WHERE room = $4 AND banned  = characters.id) AND
				NOT EXISTS (SELECT * FROM rooms_invited_characters WHERE room = $4 AND invited = characters.id) AND
				NOT EXISTS (SELECT * FROM dismiss_list WHERE id = dismiss)
			LIMIT
				20;
		`, characterId, number, searchText, roomId)
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
