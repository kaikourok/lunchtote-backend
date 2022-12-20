package character

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

func (db *CharacterRepository) DeleteLayerItems(characterId int, itemIds *[]int) error {
	query := `
		SELECT
			EXISTS (
				SELECT
					*
				FROM
					characters_icon_layer_items
				JOIN
					characters_icon_layer_groups ON (characters_icon_layer_items.layer_group = characters_icon_layer_groups.id)
				JOIN
					characters_icon_layering_groups ON (characters_icon_layer_groups.layering_group = characters_icon_layering_groups.id)
				WHERE
					characters_icon_layer_items.id IN (:list) AND
					characters_icon_layering_groups.character != :character
			);
	`

	query, args, err := sqlx.Named(query, map[string]interface{}{
		"character": characterId,
		"list":      *itemIds,
	})
	if err != nil {
		return err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}

	row := db.QueryRowx(db.Rebind(query), args...)

	var existsOtherCharacterItems bool
	err = row.Scan(&existsOtherCharacterItems)
	if err != nil {
		return err
	}
	if existsOtherCharacterItems {
		return errors.New("自身のものではないレイヤーアイテムが含まれています")
	}

	query = `
		DELETE FROM
			characters_icon_layer_items
		WHERE
			id IN (?);
	`

	query, args, err = sqlx.In(query, *itemIds)
	if err != nil {
		return err
	}

	_, err = db.Exec(db.Rebind(query), args...)

	return err
}
