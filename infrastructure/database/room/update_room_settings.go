package room

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *RoomRepository) UpdateRoomSettings(characterId, roomId int, room *model.Room) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		row := tx.QueryRowx(`
			SELECT
				master
			FROM
				rooms
			WHERE
				id = $1;
		`, roomId)

		var masterId int
		err := row.Scan(&masterId)
		if err != nil {
			return err
		}
		if masterId != characterId {
			return errors.New("トークルームの管理者ではありません")
		}

		_, err = tx.Exec(`
			UPDATE
				rooms
			SET
				title                = $2,
				summary              = $3,
				description          = $4,
				searchable           = $5,
				allow_recommendation = $6,
				children_referable   = $7
			WHERE
				id = $1;
		`,
			roomId,
			room.Title,
			room.Summary,
			room.Description,
			room.Searchable,
			room.AllowRecommendation,
			room.ChildrenReferable,
		)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			DELETE FROM
				rooms_tags
			WHERE
				room = $1;
		`, roomId)
		if err != nil {
			return err
		}

		for i := range room.Tags {
			_, err = tx.Exec(`
				INSERT INTO rooms_tags (
					room,
					tag
				) VALUES (
					$1,
					$2
				);
			`, roomId, room.Tags[i])
			if err != nil {
				return err
			}
		}

		return nil
	})
}
