package character

import "errors"

func (db *CharacterRepository) CreateLayerGroup(characterId, groupId int, name string) (id int, err error) {
	var master int
	row := db.QueryRowx(`
		SELECT
			character
		FROM
			characters_icon_layering_groups
		WHERE
			id = $1;
	`, groupId)

	err = row.Scan(&master)
	if err != nil {
		return 0, err
	}

	if characterId != master {
		return 0, errors.New("レイヤリンググループの所有者ではありません")
	}

	row = db.QueryRowx(`
		INSERT INTO characters_icon_layer_groups (
			layering_group,
			name,
			layer_order
		) VALUES (
			$1,
			$2,
			COALESCE(
				(
					SELECT
						layer_order + 1
					FROM 
						characters_icon_layer_groups
					WHERE
						layering_group = $1
					ORDER BY
						layer_order DESC
					LIMIT
						1
				)
			, 0)
		)
		RETURNING
			id;
	`, groupId, name)

	err = row.Scan(&id)
	return
}
