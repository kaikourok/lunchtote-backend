package character

import "errors"

func (db *CharacterRepository) DeleteLayerGroup(characterId, groupId int) error {
	var master int
	row := db.QueryRowx(`
		SELECT
			characters_icon_layering_groups.character
		FROM
			characters_icon_layer_groups
		JOIN
			characters_icon_layering_groups ON (characters_icon_layer_groups.layering_group = characters_icon_layering_groups.id)
		WHERE
			characters_icon_layer_groups.id = $1;
	`, groupId)

	err := row.Scan(&master)
	if err != nil {
		return err
	}

	if characterId != master {
		return errors.New("レイヤリンググループの所有者ではありません")
	}

	_, err = db.Exec(`
		DELETE FROM
			characters_icon_layer_groups
		WHERE
			id = $1;
	`, groupId)
	return err
}
