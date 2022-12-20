package character

import "errors"

func (db *CharacterRepository) DeleteLayeringGroup(characterId int, groupId int) error {
	var master int
	row := db.QueryRowx(`
		SELECT
			character
		FROM
			characters_icon_layering_groups
		WHERE
			id = $1;
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
			characters_icon_layering_groups
		WHERE
			id = $1;
	`, groupId)
	return err
}
