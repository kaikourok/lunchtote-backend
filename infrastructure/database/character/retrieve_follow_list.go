package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/lib/pq"
)

func (db *CharacterRepository) RetrieveFollowList(userId, targetId int) (list *[]model.CharacterListItem, err error) {
	rows, err := db.Queryx(`
		WITH
			follow_list   AS (SELECT followed FROM follows WHERE follower = $2),
			follower_list AS (SELECT follower FROM follows WHERE followed = $2),
			mute_list     AS (SELECT muted    FROM mutes   WHERE muter    = $2),
			block_list    AS (SELECT blocked  FROM blocks  WHERE blocker  = $2)
		SELECT
			characters.id,
			characters.name,
			characters.summary,
			characters.Mainicon,
			ARRAY_REMOVE(ARRAY_AGG(characters_tags.tag ORDER BY characters_tags.id), NULL),
			characters.id IN follow_list,
			characters.id IN follower_list,
			characters.id IN mute_list,
			characters.id IN block_list
		FROM
			follows
		JOIN
			characters ON (follows.followed = characters.id AND follows.follower = $1)
		LEFT JOIN
			characters_tags ON (characters.id = characters_tags.character)
		WHERE
			characters.deleted_at IS NULL AND
			characters.administrator = false
		GROUP BY
			characters.id,
			characters.name,
			characters.summary,
			characters.Mainicon
		ORDER BY
			characters_mutes.muted_at DESC;
	`, targetId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	characterList := make([]model.CharacterListItem, 0, 64)
	for i := 0; rows.Next(); i++ {
		var listItem model.CharacterListItem

		err = rows.Scan(
			&listItem.Id,
			&listItem.Name,
			&listItem.Summary,
			&listItem.Mainicon,
			pq.Array(&listItem.Tags),
			&listItem.IsFollowing,
			&listItem.IsFollowed,
			&listItem.IsMuting,
			&listItem.IsBlocking,
		)
		if err != nil {
			return nil, err
		}

		characterList = append(characterList, listItem)
	}

	return &characterList, nil
}
