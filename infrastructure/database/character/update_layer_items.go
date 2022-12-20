package character

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) UpdateLayerItems(characterId, groupId int, items *[]model.CharacterIconLayerItemEditData) (result *[]model.CharacterIconLayerItem, err error) {
	resultItems := make([]model.CharacterIconLayerItem, 0, len(*items))

	err = db.ExecTx(func(tx *sqlx.Tx) error {
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

		_, err = tx.Exec(`
			DELETE FROM
				characters_icon_layer_items 
			WHERE
				layer_group = $1;
		`, groupId)
		if err != nil {
			return err
		}

		inserts := make([]struct {
			LayerGroup int    `db:"layer_group"`
			Path       string `db:"path"`
		}, len(*items))
		for i, item := range *items {
			inserts[i].LayerGroup = groupId
			inserts[i].Path = item.Path
		}

		_, err = tx.NamedExec(`
			INSERT INTO characters_icon_layer_items (
				layer_group,
				path
			) VALUES (
				:layer_group,
				:path
			)
		`, inserts)
		if err != nil {
			return err
		}

		rows, err := tx.Queryx(`
			SELECT
				id,
				path
			FROM
				characters_icon_layer_items
			WHERE
				layer_group = $1
			ORDER BY
				id;
		`, groupId)
		if err != nil {
			return err
		}

		for rows.Next() {
			var result model.CharacterIconLayerItem
			err = rows.Scan(
				&result.Id,
				&result.Path,
			)
			if err != nil {
				return err
			}

			resultItems = append(resultItems, result)
		}

		return nil
	})

	return &resultItems, err
}
