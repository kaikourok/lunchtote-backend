package room

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *RoomRepository) UpdateRoomMessageFetchConfigOrders(characterId int, orders *[]model.RoomMessageFetchConfigOrder) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		rows, err := tx.Queryx(`
			SELECT
				id
			FROM
				message_fetch_configs
			WHERE
				master = $1;
		`, characterId)
		if err != nil {
			return err
		}
		defer rows.Close()

		configIds := make([]int, 0, len(*orders))
		for rows.Next() {
			var configId int
			err = rows.Scan(&configId)
			if err != nil {
				return err
			}

			configIds = append(configIds, configId)
		}

		if len(configIds) != len(*orders) {
			return errors.New("指定数が不足しています")
		}

		for _, order := range *orders {
			_, err = tx.Exec(`
				UPDATE
				  message_fetch_configs
				SET
					config_order = $2
				WHERE
					id = $1;
			`, order.Config, order.Order)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
