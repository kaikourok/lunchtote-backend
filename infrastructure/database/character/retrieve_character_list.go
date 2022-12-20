package character

import (
	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/lib/pq"
)

func (db *CharacterRepository) RetrieveCharacterList(id *int, start, end int) (list *[]model.CharacterListItem, maxId int, err error) {
	var rows *sqlx.Rows
	if id == nil {
		// 非ログインユーザー視点
		rows, err = db.Queryx(`
			WITH max_id AS (
				SELECT
					id
				FROM
					characters			
				WHERE
					deleted_at IS NULL AND
					administrator = false
				ORDER BY
					id DESC
				LIMIT
					1
			)
			SELECT
				characters.id,
				characters.name,
				characters.nickname,
				characters.summary,
				characters.Mainicon,
				ARRAY_REMOVE(ARRAY_AGG(characters_tags.tag ORDER BY characters_tags.id), NULL),
				COALESCE((SELECT id FROM max_id), 0) AS max_id,
				null,
				null,
				null
			FROM
				characters
			LEFT JOIN
				characters_tags ON (characters.id = characters_tags.character)
			WHERE
				characters.id BETWEEN $1 AND $2 AND
				characters.deleted_at IS NULL AND
				characters.administrator = false
			GROUP BY
				characters.id,
				characters.name,
				characters.nickname,
				characters.summary,
				characters.Mainicon,
				max_id
			ORDER BY
				characters.id;
		`, start, end)
	} else {
		// ログインユーザー視点
		rows, err = db.Queryx(`
			WITH
				follow_list   AS (SELECT followed FROM follows WHERE follower = $3),
				follower_list AS (SELECT follower FROM follows WHERE followed = $3),
				mute_list     AS (SELECT muted    FROM mutes   WHERE muter    = $3),
				max_id AS (
					SELECT
						id
					FROM
						characters			
					WHERE
						deleted_at IS NULL AND
						administrator = false
					ORDER BY
						id DESC
					LIMIT
						1
				)
			SELECT
				characters.id,
				characters.name,
				characters.nickname,
				characters.summary,
				characters.Mainicon,
				ARRAY_REMOVE(ARRAY_AGG(characters_tags.tag ORDER BY characters_tags.id), NULL),
				COALESCE((SELECT id FROM max_id), 0) AS max_id,
				characters.id IN (SELECT * FROM follow_list),
				characters.id IN (SELECT * FROM follower_list),
				characters.id IN (SELECT * FROM mute_list)
			FROM
				characters
			LEFT JOIN
				characters_tags ON (characters.id = characters_tags.character)
			WHERE
				characters.id BETWEEN $1 AND $2 AND
				characters.deleted_at IS NULL AND
				characters.administrator = false AND
				characters.id NOT IN (
					SELECT
						blocked
					FROM
						blocks
					WHERE
						blocker = $3
				) AND 
				characters.id NOT IN (
					SELECT
						blocker
					FROM
						blocks
					WHERE
						blocked = $3
				)
			GROUP BY
				characters.id,
				characters.name,
				characters.nickname,
				characters.summary,
				characters.Mainicon,
				max_id
			ORDER BY
				characters.id;
		`, start, end, *id)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	characterList := make([]model.CharacterListItem, 0, end-start+1)

	for i := 0; rows.Next(); i++ {
		var listItem model.CharacterListItem

		err = rows.Scan(
			&listItem.Id,
			&listItem.Name,
			&listItem.Nickname,
			&listItem.Summary,
			&listItem.Mainicon,
			pq.Array(&listItem.Tags),
			&maxId,
			&listItem.IsFollowing,
			&listItem.IsFollowed,
			&listItem.IsMuting,
		)
		if err != nil {
			return nil, 0, err
		}

		characterList = append(characterList, listItem)
	}

	return &characterList, maxId, nil
}
