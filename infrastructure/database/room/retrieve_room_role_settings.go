package room

import (
	"encoding/json"

	model "github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *RoomRepository) RetrieveRoomRoleSettings(roomId int) (roles *[]model.RoomRoleWithMembers, master int, title string, err error) {

	rows, err := db.Queryx(`
		SELECT
			rooms_roles.id,
			rooms.master,
			rooms.title,
			rooms_roles.priority,
			COALESCE(rooms_roles.name, ''),
			rooms_roles.write,
			rooms_roles.ban,
			rooms_roles.invite,
			rooms_roles.use_reply,
			rooms_roles.use_secret,
			rooms_roles.delete_other_message,
			rooms_roles.type,
			JSON_AGG(JSON_BUILD_OBJECT(
				'id',       characters.id,
				'name',     characters.name,
				'mainicon', characters.mainicon
			))
		FROM
			rooms
		JOIN
			rooms_roles ON (rooms.id = $1 AND rooms.id = rooms_roles.room)
		LEFT JOIN
			rooms_roles_members ON (rooms_roles.id = rooms_roles_members.role)
		LEFT JOIN
			characters ON (rooms_roles_members.character = characters.id)
		GROUP BY
			rooms_roles.id,
			rooms.master,
			rooms.title,
			rooms_roles.priority,
			rooms_roles.name,
			rooms_roles.write,
			rooms_roles.ban,
			rooms_roles.invite,
			rooms_roles.use_reply,
			rooms_roles.use_secret,
			rooms_roles.delete_other_message,
			rooms_roles.type
		ORDER BY
			rooms_roles.priority DESC;
	`, roomId)

	if err != nil {
		return nil, 0, "", err
	}
	defer rows.Close()

	rolesSlice := make([]model.RoomRoleWithMembers, 0, 32)

	for rows.Next() {
		var role model.RoomRoleWithMembers
		var roleMembersJson string
		err = rows.Scan(
			&role.Id,
			&master,
			&title,
			&role.Priority,
			&role.Name,
			&role.Write,
			&role.Ban,
			&role.Invite,
			&role.UseReply,
			&role.UseSecret,
			&role.DeleteOtherMessage,
			&role.Type,
			&roleMembersJson,
		)
		if err != nil {
			return nil, 0, "", err
		}

		err = json.Unmarshal([]byte(roleMembersJson), &role.Members)
		if err != nil {
			return nil, 0, "", err
		}

		rolesSlice = append(rolesSlice, role)
	}

	return &rolesSlice, master, title, nil
}
