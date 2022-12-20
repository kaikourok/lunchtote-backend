package room

import (
	"encoding/json"
	"errors"

	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *RoomRepository) RetrieveRoomMembers(userId, roomId int) (members *[]model.RomeMemberWithRoles, err error) {
	row := db.QueryRowx(`
		SELECT
			EXISTS (
				SELECT
					*
				FROM
					rooms
				WHERE
					id     = $1 AND
					master = $2
			);
	`, roomId, userId)

	var isMaster bool
	err = row.Scan(&isMaster)
	if err != nil {
		return nil, err
	}

	if !isMaster {
		return nil, errors.New("管理権限がありません")
	}

	rows, err := db.Queryx(`
		SELECT
			characters.id,
			characters.name,
			characters.mainicon
			JSON_AGG(JSON_BUILD_OBJECT(
				'id',   rooms_roles.id,
				'name', rooms_roles_members.
			))
		FROM
			rooms_members
		JOIN
			characters ON (rooms_members.character = characters.id AND rooms_members.room = $1)
		LEFT JOIN
			rooms_roles_members ON (rooms_members.character = rooms_roles_members.character)
		LEFT JOIN
			rooms_roles ON (rooms_roles_members.role = rooms_roles)
		WHERE
			rooms_members.room = $1
		ORDER BY
			rooms_members.id;
	`, roomId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	membersSlice := make([]model.RomeMemberWithRoles, 0, 64)

	for rows.Next() {
		var member model.RomeMemberWithRoles
		var memberRolesJson string
		err = rows.Scan(
			&member.Id,
			&member.Name,
			&member.Mainicon,
			&memberRolesJson,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(memberRolesJson), &member.Roles)
		if err != nil {
			return nil, err
		}

		membersSlice = append(membersSlice, member)
	}

	return &membersSlice, nil
}
