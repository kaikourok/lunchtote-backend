package character

import (
	"encoding/json"

	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/lib/pq"
)

func (db *CharacterRepository) RetrieveProfile(userId *int, targetId int) (*model.Profile, error) {
	row := db.QueryRowx(`
		SELECT
			id,
			name,
			nickname,
			summary,
			profile,
			ARRAY(SELECT path FROM characters_profile_images WHERE character = $1),
			ARRAY(SELECT tag  FROM characters_tags           WHERE character = $1),
			(SELECT count(*) FROM follows WHERE follower = $1),
			(SELECT count(*) FROM follows WHERE followed = $1),
			COALESCE(
				(
					SELECT 
						json_agg(row_to_json(icon_rows))
					FROM
						(
							SELECT
								path
							FROM
								characters_icons
							WHERE
								character = $1
							ORDER BY
								id
							LIMIT
								30
						) icon_rows
				)
			, '[]'::JSON) AS icons_json,
			COALESCE(
				(
					SELECT 
						json_agg(row_to_json(character_diaries))
					FROM
						(
							SELECT
								nth,
								title
							FROM
								diaries
							WHERE
								character = $1
							ORDER BY
								nth
						) character_diaries
				)
			, '[]'::JSON) AS diaries_json,
			COALESCE((SELECT true FROM follows WHERE follower = $2 AND followed = $1), false),
			COALESCE((SELECT true FROM follows WHERE followed = $2 AND follower = $1), false),
			COALESCE((SELECT true FROM mutes   WHERE muter    = $2 AND muted    = $1), false),
			COALESCE((SELECT true FROM blocks  WHERE blocker  = $2 AND blocked  = $1), false),
			COALESCE((SELECT true FROM blocks  WHERE blocked  = $2 AND blocker  = $1), false)
		FROM
			characters
		WHERE
			id = $1 AND
			deleted_at IS NULL AND
			administrator = false;
	`, targetId, userId)

	character := model.Profile{}

	var iconsJsonReader, diariesJsonReader string
	err := row.Scan(
		&character.Id,
		&character.Name,
		&character.Nickname,
		&character.Summary,
		&character.Profile,
		pq.Array(&character.ProfileImages),
		pq.Array(&character.Tags),
		&character.FollowingNumber,
		&character.FollowedNumber,
		&iconsJsonReader,
		&diariesJsonReader,
		&character.IsFollowing,
		&character.IsFollowed,
		&character.IsMuting,
		&character.IsBlocking,
		&character.IsBlocked,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(iconsJsonReader), &character.Icons)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(diariesJsonReader), &character.ExistingDiaries)
	if err != nil {
		return nil, err
	}

	return &character, nil
}
