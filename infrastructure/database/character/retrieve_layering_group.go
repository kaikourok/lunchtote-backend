package character

import (
	"encoding/json"
	"errors"

	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) RetrieveLayeringGroup(characterId, groupId int) (layeringGroup *model.CharacterIconLayeringGroup, err error) {
	layeringGroup = &model.CharacterIconLayeringGroup{}
	layeringGroup.Id = groupId

	var master int
	row := db.QueryRowx(`
		SELECT
			character,
			name
		FROM
			characters_icon_layering_groups
		WHERE
			id = $1;
	`, groupId)

	err = row.Scan(&master, &layeringGroup.Name)
	if err != nil {
		return nil, err
	}

	if characterId != master {
		return nil, errors.New("レイヤリンググループの所有者ではありません")
	}

	rows, err := db.Queryx(`
		SELECT
			characters_icon_layer_groups.id,
			characters_icon_layer_groups.name,
			COALESCE(
				JSON_AGG(JSON_BUILD_OBJECT(
					'id',   characters_icon_layer_items.id,
					'path', characters_icon_layer_items.path
				)) FILTER (WHERE characters_icon_layer_items.path IS NOT NULL),
				'[]'
			)
		FROM
			characters_icon_layer_groups
		LEFT JOIN
			characters_icon_layer_items ON (characters_icon_layer_groups.id = characters_icon_layer_items.layer_group)
		WHERE
			characters_icon_layer_groups.layering_group = $1
		GROUP BY
			characters_icon_layer_groups.id,
			characters_icon_layer_groups.name
		ORDER BY
			characters_icon_layer_groups.layer_order DESC;
	`, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	layeringGroup.Layers = make([]model.CharacterIconLayerGroup, 0, 64)
	for rows.Next() {
		var layerGroup model.CharacterIconLayerGroup
		var layerGroupItemsJson string
		err = rows.Scan(
			&layerGroup.Id,
			&layerGroup.Name,
			&layerGroupItemsJson,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(layerGroupItemsJson), &layerGroup.Items)
		if err != nil {
			return nil, err
		}

		layeringGroup.Layers = append(layeringGroup.Layers, layerGroup)
	}

	rows, err = db.Queryx(`
		SELECT
			id,
			name,
			x,
			y,
			scale,
			rotate
		FROM
			characters_icon_process_schemas
		WHERE
			layering_group = $1
		ORDER BY
			id;
	`, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	layeringGroup.Processes = make([]model.CharacterIconProcessSchemaEditData, 0, 64)
	for rows.Next() {
		var process model.CharacterIconProcessSchemaEditData
		err = rows.Scan(
			&process.Id,
			&process.Name,
			&process.X,
			&process.Y,
			&process.Scale,
			&process.Rotate,
		)
		if err != nil {
			return nil, err
		}

		layeringGroup.Processes = append(layeringGroup.Processes, process)
	}

	return
}
