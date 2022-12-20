package character

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/library/slice"
)

func (db *CharacterRepository) UpdateLayerGroupOrders(characterId, groupId int, orders *[]model.CharacterIconLayerGroupOrder) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		rows, err := tx.Queryx(`
			SELECT
				characters_icon_layering_groups.character,
				characters_icon_layer_groups.id
			FROM
				characters_icon_layer_groups
			JOIN
				characters_icon_layering_groups ON (characters_icon_layering_groups.id = characters_icon_layer_groups.layering_group)
			WHERE
				characters_icon_layering_groups.id = $1;
		`, groupId)
		if err != nil {
			return err
		}
		defer rows.Close()

		layerIds := make([]int, 0, len(*orders))
		for rows.Next() {
			var masterId, layerId int
			err = rows.Scan(&masterId, &layerId)
			if err != nil {
				return err
			}

			if masterId != characterId {
				return errors.New("レイヤー順を変更する権限がありません")
			}

			layerIds = append(layerIds, layerId)
		}

		if len(layerIds) != len(*orders) {
			return errors.New("指定数が不足しています")
		}

		for _, order := range *orders {
			if !slice.Contains(order.LayerGroup, &layerIds) {
				return errors.New("レイヤグループ外のIDが指定されています")
			}
		}

		for _, order := range *orders {
			_, err = tx.Exec(`
				UPDATE
					characters_icon_layer_groups
				SET
					layer_order = $2
				WHERE
					id = $1;
			`, order.LayerGroup, order.Order)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
