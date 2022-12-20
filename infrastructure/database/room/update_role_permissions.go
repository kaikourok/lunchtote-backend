package room

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *RoomRepository) UpdateRolePermissions(characterId, roleId int, roleName string, role *model.RoomRolePermission) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		row := tx.QueryRowx(`
			SELECT
				rooms.id,
				rooms.master,
				rooms_roles.type
			FROM
				rooms_roles
			JOIN
				rooms ON (rooms_roles.room = rooms.id) 
			WHERE
				rooms_roles.id = $1;
		`, roleId)

		var roomId, masterId int
		var roleType string
		err := row.Scan(&roomId, &masterId, &roleType)
		if err != nil {
			return err
		}

		if masterId != characterId {
			return errors.New("権限を変更する権限がありません")
		}

		err = db.roleTablesAdvisoryLock(roomId, tx)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE
				rooms_roles
			SET
				name                 = $2,
				write                = $3,
				ban                  = $4,
				invite               = $5,
				use_reply            = $6,
				use_secret           = $7,
				delete_other_message = $8,
				create_children_room = $9
			WHERE
				rooms_roles.id = $1;
		`,
			roleId,
			roleName,
			role.Write,
			role.Ban,
			role.Invite,
			role.UseReply,
			role.UseSecret,
			role.DeleteOtherMessage,
			role.CreateChildrenRoom,
		)

		if err != nil {
			return err
		}

		return db.attachPermissions(roomId, tx)
	})
}
