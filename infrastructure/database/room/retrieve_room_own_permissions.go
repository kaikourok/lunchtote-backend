package room

import (
	"database/sql"

	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *RoomRepository) RetrieveRoomOwnPermissions(characterId int, roomId int) (permissions *model.RoomMemberPermission, banned bool, err error) {
	permissions = &model.RoomMemberPermission{}

	row := db.QueryRowx(`
		SELECT
			write,
			ban,
			invite,
			use_reply,
			use_secret,
			delete_other_message,
			create_children_room
		FROM
			rooms_members
		WHERE
			room = $1 AND member = $2
	`, roomId, characterId)

	err = row.Scan(
		&permissions.Write,
		&permissions.Ban,
		&permissions.Invite,
		&permissions.UseReply,
		&permissions.UseSecret,
		&permissions.DeleteOtherMessage,
		&permissions.CreateChildrenRoom,
	)
	if err == nil {
		return permissions, false, nil
	} else if err != sql.ErrNoRows {
		return nil, false, err
	}

	row = db.QueryRowx(`
		SELECT
			EXISTS (
				SELECT
					*
				FROM
					rooms_invited_characters
				WHERE
					room = $1 AND invited = $2
			),
			EXISTS (
				SELECT
					*
				FROM
					rooms_banned_characters
				WHERE
					room = $1 AND banned = $2
			);
	`, roomId, characterId)

	var isInvited, isBanned bool
	err = row.Scan(&isInvited, &isBanned)
	if err != nil {
		return nil, false, err
	}
	if isBanned {
		return nil, true, err
	}

	var roleType string
	if isInvited {
		roleType = "INVITED"
	} else {
		roleType = "VISITOR"
	}

	row = db.QueryRowx(`
		SELECT
			write,
			ban,
			invite,
			use_reply,
			use_secret,
			delete_other_message,
			create_children_room
		FROM
			rooms_roles
		WHERE
			room = $1 AND type = $2
	`, roomId, roleType)
	err = row.Scan(
		&permissions.Write,
		&permissions.Ban,
		&permissions.Invite,
		&permissions.UseReply,
		&permissions.UseSecret,
		&permissions.DeleteOtherMessage,
		&permissions.CreateChildrenRoom,
	)
	if err != nil {
		return nil, false, err
	}

	return permissions, false, nil
}
