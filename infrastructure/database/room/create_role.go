package room

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *RoomRepository) CreateRole(characterId, roomId int, roleName string, role *model.RoomRolePermission) (roleId int, err error) {
	row := db.QueryRowx(`
		SELECT
			rooms.master
		FROM
			rooms
		WHERE
			rooms.id = $1;
	`, roomId)

	var masterId int
	err = row.Scan(&masterId)
	if err != nil {
		return 0, err
	}

	if masterId != characterId {
		return 0, errors.New("ロールを作成する権限がありません")
	}

	err = db.ExecTx(func(tx *sqlx.Tx) error {
		if err != nil {
			return err
		}

		row = tx.QueryRowx(`
			INSERT INTO rooms_roles (
				room,
				priority,
				name,
				write,
				ban,
				invite,
				use_reply,
				use_secret,
				delete_other_message,
				create_children_room
			) VALUES (
				$1,
				(SELECT priority + 1 FROM rooms_roles WHERE room = $1 ORDER BY priority DESC LIMIT 1),
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
				$9
			)
			RETURNING
				rooms_roles.id;
		`,
			roomId,
			roleName,
			role.Write,
			role.Ban,
			role.Invite,
			role.UseReply,
			role.UseSecret,
			role.DeleteOtherMessage,
			role.CreateChildrenRoom,
		)

		err = row.Scan(&roleId)
		if err != nil {
			return err
		}

		err = db.attachPermissions(roomId, tx)
		if err != nil {
			return err
		}

		return nil
	})

	return roleId, err
}
