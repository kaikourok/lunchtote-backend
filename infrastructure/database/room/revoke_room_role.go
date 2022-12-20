package room

import "errors"

func (db *RoomRepository) RevokeRoomRole(userId, targetId, role int) error {
	row := db.QueryRowx(`
		SELECT
			rooms.master
		FROM
			rooms_roles
		JOIN
			rooms ON (rooms_roles.room = rooms.id) 
		WHERE
			rooms_roles.id = $1;
	`, role)

	var masterId int
	err := row.Scan(&masterId)
	if err != nil {
		return err
	}

	if masterId != userId {
		return errors.New("ロールを剥奪する権限がありません")
	}

	_, err = db.Exec(`
		DELETE FROM
			rooms_roles_members
		WHERE
			role      = $2 AND
			character = $1;
	`, targetId, role)
	if err != nil {
		return err
	}

	return nil
}
