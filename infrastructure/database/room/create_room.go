package room

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *RoomRepository) CreateRoom(characterId int, room *model.Room) (roomId int, err error) {
	err = db.ExecTx(func(tx *sqlx.Tx) error {
		if room.ParentRoom != nil {
			row := tx.QueryRowx(`
				SELECT
					EXISTS (
						SELECT
							*
						FROM
							rooms_members
						WHERE
							room   = $1 AND
							member = $2 AND
							create_children_room = true
					);
			`, room.ParentRoom, characterId)

			var belongable bool
			err := row.Scan(&belongable)
			if err != nil {
				return err
			}
			if !belongable {
				return errors.New("所属ルームを作成する権限がありません")
			}
		}

		row := tx.QueryRowx(`
			INSERT INTO rooms (
				master,
				title,
				summary,
				description,
				searchable,
				allow_recommendation,
				children_referable,
				belong
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8
			)
			RETURNING
				seq;
		`,
			characterId,
			room.Title,
			room.Summary,
			room.Description,
			room.Searchable,
			room.AllowRecommendation,
			room.ChildrenReferable,
			room.ParentRoom,
		)

		var seq int
		err := row.Scan(&seq)
		if err != nil {
			return err
		}

		row = tx.QueryRowx(`
			UPDATE rooms SET
				id = (SELECT count(*) FROM rooms WHERE official = false)
			WHERE
				seq = $1
			RETURNING
				id;
		`, seq)

		err = row.Scan(&roomId)
		if err != nil {
			return err
		}

		type InsertRoomTagStruct struct {
			Room int    `db:"room"`
			Tag  string `db:"tag"`
		}

		if 0 < len(room.Tags) {
			inserts := make([]InsertRoomTagStruct, len(room.Tags))
			for i, tag := range room.Tags {
				inserts[i].Room = roomId
				inserts[i].Tag = tag
			}

			_, err = tx.NamedExec(`
			INSERT INTO rooms_tags (room, tag) VALUES (:room, :tag)
		`, inserts)
			if err != nil {
				return err
			}
		}

		_, err = tx.Exec(`
			INSERT INTO rooms_roles (
				room,
				priority,
				read,
				write,
				ban,
				invite,
				use_reply,
				use_secret,
				delete_other_message,
				create_children_room,
				type
			) VALUES (
				$1,
				-4,
				true,
				true,
				false,
				false,
				true,
				true,
				false,
				false,
				'VISITOR'
			), (
				$1,
				-3,
				true,
				true,
				false,
				false,
				true,
				true,
				false,
				false,
				'INVITED'
			), (
				$1,
				-2,
				true,
				true,
				false,
				false,
				true,
				true,
				false,
				false,
				'DEFAULT'
			), (
				$1,
				-1,
				true,
				true,
				true,
				true,
				true,
				true,
				true,
				true,
				'MASTER'
			);
		`, roomId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO rooms_members (
				room,
				member,
				write,
				ban,
				invite,
				use_reply,
				use_secret,
				delete_other_message,
				create_children_room,
				type
			) VALUES (
				$1,
				$2,
				true,
				true,
				true,
				true,
				true,
				true,
				true,
				'MASTER'
			);
		`, roomId, characterId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO rooms_roles_members (
				role,
				character
			) VALUES (
				(SELECT id FROM rooms_roles WHERE room = $1 AND type = 'MASTER'),
				$2
			);
		`, roomId, characterId)
		if err != nil {
			return err
		}

		return nil
	})

	return
}
