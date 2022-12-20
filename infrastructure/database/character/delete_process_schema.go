package character

import "errors"

func (db *CharacterRepository) DeleteProcessSchema(characterId, processId int) error {
	var master int
	row := db.QueryRowx(`
		SELECT
			characters_icon_layering_groups.character
		FROM
			characters_icon_process_schemas
		JOIN
			characters_icon_layering_groups ON (characters_icon_process_schemas.layering_group = characters_icon_layering_groups.id)
		WHERE
			characters_icon_process_schemas.id = $1;
	`, processId)

	err := row.Scan(&master)
	if err != nil {
		return err
	}

	if characterId != master {
		return errors.New("レイヤリンググループの所有者ではありません")
	}

	_, err = db.Exec(`
		DELETE FROM 
			characters_icon_process_schemas
		WHERE
			id = $1;
	`, processId)

	return err
}
