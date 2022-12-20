package room

import (
	"errors"
)

func (db *RoomRepository) GrantRoomRole(userId, targetId, roleId int) error {
	row := db.QueryRowx(`
		SELECT
			rooms.master
		FROM
			rooms_roles
		JOIN
			rooms ON (rooms_roles.room = rooms.id) 
		WHERE
			rooms_roles.id = $1;
	`, roleId)

	var masterId int
	err := row.Scan(&masterId)
	if err != nil {
		return err
	}

	if masterId != userId {
		return errors.New("ロールを付与する権限がありません")
	}

	_, err = db.Exec(`
		INSERT INTO rooms_roles_members (
			role,
			character
		) VALUES (
			$2,
			$1
		);
	`, targetId, roleId)
	if err != nil {
		return err
	}

	return nil
}
