package room

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

func (db *RoomRepository) DeleteRole(characterId, roleId int) error {
	row := db.QueryRowx(`
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
		return errors.New("ロールを削除する権限がありません")
	}

	if roleType != "MEMBER" {
		return errors.New("このロールは削除できません")
	}

	return db.ExecTx(func(tx *sqlx.Tx) error {
		err = db.roleTablesAdvisoryLock(roomId, tx)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			DELETE FROM
				rooms_roles_members
			WHERE
				role = $1
		`, roleId)

		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			DELETE FROM
				rooms_roles
			WHERE
				id = $1
		`, roleId)

		if err != nil {
			return err
		}

		return db.attachPermissions(roomId, tx)
	})
}
