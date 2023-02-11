package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/lib/pq"
)

func (db *CharacterRepository) RetrieveBlockList(id int) (list *[]model.GeneralCharacterListItem, err error) {
	rows, err := db.Queryx(`
		WITH
			mute_list AS (SELECT followed FROM follows WHERE follower = $1)
		SELECT
			characters.id,
			characters.name,
			characters.summary,
			characters.mainicon,
			ARRAY_REMOVE(ARRAY_AGG(characters_tags.tag ORDER BY characters_tags.id), NULL),
			characters.id IN (SELECT * FROM mute_list)
		FROM
			blocks
		JOIN
			characters ON (blocks.blocked = characters.id AND blocks.blocker = $1)
		LEFT JOIN
			characters_tags ON (characters.id = characters_tags.character)
		WHERE
			characters.deleted_at IS NULL AND
			characters.administrator = false
		GROUP BY
			characters.id,
			characters.name,
			characters.summary,
			characters.mainicon,
			blocks.blocked_at
		ORDER BY
			blocks.blocked_at DESC;
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	characterList := make([]model.GeneralCharacterListItem, 0, 64)
	for i := 0; rows.Next(); i++ {
		var listItem model.GeneralCharacterListItem

		err = rows.Scan(
			&listItem.Id,
			&listItem.Name,
			&listItem.Summary,
			&listItem.Mainicon,
			pq.Array(&listItem.Tags),
			&listItem.IsMuting,
		)
		if err != nil {
			return nil, err
		}

		characterList = append(characterList, listItem)
	}

	return &characterList, nil
}
