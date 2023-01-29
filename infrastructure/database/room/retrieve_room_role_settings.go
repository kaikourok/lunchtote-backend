package room

import (
	model "github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *RoomRepository) RetrieveRoomRoleSettings(roomId int) (roles []model.RoomRole, master int, err error) {
	rows, err := db.Queryx(`
		SELECT
			rooms_roles.id,
			rooms.master,
			rooms_roles.priority,
			COALESCE(rooms_roles.name, ''),
			rooms_roles.write,
			rooms_roles.ban,
			rooms_roles.invite,
			rooms_roles.use_reply,
			rooms_roles.use_secret,
			rooms_roles.delete_other_message,
			rooms_roles.type
		FROM
			rooms
		JOIN
			rooms_roles ON (rooms.id = $1 AND rooms.id = rooms_roles.room)
		ORDER BY
			rooms_roles.priority DESC;
	`, roomId)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	rolesSlice := make([]model.RoomRole, 0, 32)

	for rows.Next() {
		var role model.RoomRole
		err = rows.Scan(
			&role.Id,
			&master,
			&role.Priority,
			&role.Name,
			&role.Write,
			&role.Ban,
			&role.Invite,
			&role.UseReply,
			&role.UseSecret,
			&role.DeleteOtherMessage,
			&role.Type,
		)
		if err != nil {
			return nil, 0, err
		}

		rolesSlice = append(rolesSlice, role)
	}

	return rolesSlice, master, nil
}
