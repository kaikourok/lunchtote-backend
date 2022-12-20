package room

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

func (db *RoomRepository) roleTablesAdvisoryLock(roomId int, tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		UPDATE
			rooms_role_tables_advisory_locker
		SET
			role_version = role_version + 1
		WHERE
			room = $1;
	`, roomId)

	if err != nil {
		return err
	}

	return nil
}

func (db *RoomRepository) attachPermissions(roomId int, tx *sqlx.Tx) error {
	type role struct {
		id                 int
		roleType           string
		read               *bool
		write              *bool
		ban                *bool
		invite             *bool
		useReply           *bool
		useSecret          *bool
		deleteOtherMessage *bool
	}

	rows, err := tx.Queryx(`
		SELECT
			id,
			type,
			read,
			write,
			ban,
			invite,
			use_reply,
			use_secret,
			delete_other_message
		FROM
			rooms_roles
		WHERE
			rooms_roles.room = $1
		ORDER BY
			priority DESC;
	`, roomId)

	if err != nil {
		return err
	}
	defer rows.Close()

	roles := make([]role, 0, 64)

	for rows.Next() {
		var role role
		err = rows.Scan(
			&role.id,
			&role.roleType,
			&role.read,
			&role.write,
			&role.ban,
			&role.invite,
			&role.useReply,
			&role.useSecret,
			&role.deleteOtherMessage,
		)

		if err != nil {
			return err
		}

		roles = append(roles, role)
	}

	rows, err = tx.Queryx(`
		SELECT
			ordered.character,
			JSON_AGG(ordered.role)
		FROM
			(
				SELECT
					rooms_roles_members.character,
					rooms_roles_members.role
				FROM
					rooms_roles
				JOIN
					rooms_roles_members ON (rooms_roles.id = rooms_roles_members.role AND rooms_roles.room = $1)
				ORDER BY
					rooms_roles.priority,
					rooms_roles_members.character
			) AS ordered
		GROUP BY
			ordered.character;
	`, roomId)

	if err != nil {
		return err
	}
	defer rows.Close()

	getRoleFromId := func(roleId int) *role {
		for _, role := range roles {
			if role.id == roleId {
				return &role
			}
		}

		panic("指定IDのロールが見つかりません")
	}

	applyGreaterRolePermission := func(oldPermission bool, newPermission *bool) bool {
		if newPermission == nil {
			return oldPermission
		} else {
			return *newPermission
		}
	}

	type MemberPermissions struct {
		Room               int    `db:"room"`
		Member             int    `db:"member"`
		Read               bool   `db:"read"`
		Write              bool   `db:"write"`
		Ban                bool   `db:"ban"`
		Invite             bool   `db:"invite"`
		UseReply           bool   `db:"use_reply"`
		UseSecret          bool   `db:"use_secret"`
		DeleteOtherMessage bool   `db:"delete_other_message"`
		Type               string `db:"type"`
	}

	members := make([]MemberPermissions, 0, 64)

	for rows.Next() {
		var character int
		var rolesJsonReader string
		err = rows.Scan(&character, &rolesJsonReader)

		if err != nil {
			return err
		}

		var roles []int

		err = json.Unmarshal([]byte(rolesJsonReader), &roles)
		if err != nil {
			return err
		}

		permissions := MemberPermissions{
			Room:               roomId,
			Member:             character,
			Read:               false,
			Write:              false,
			Ban:                false,
			Invite:             false,
			UseReply:           false,
			UseSecret:          false,
			DeleteOtherMessage: false,
			Type:               "MEMBER",
		}

		for _, roleId := range roles {
			role := getRoleFromId(roleId)

			permissions.Read = applyGreaterRolePermission(permissions.Read, role.read)
			permissions.Write = applyGreaterRolePermission(permissions.Write, role.write)
			permissions.Ban = applyGreaterRolePermission(permissions.Ban, role.ban)
			permissions.Invite = applyGreaterRolePermission(permissions.Invite, role.invite)
			permissions.UseReply = applyGreaterRolePermission(permissions.UseReply, role.useReply)
			permissions.UseSecret = applyGreaterRolePermission(permissions.UseSecret, role.useSecret)
			permissions.DeleteOtherMessage = applyGreaterRolePermission(permissions.DeleteOtherMessage, role.deleteOtherMessage)

			if role.roleType == "MASTER" {
				permissions.Type = "MASTER"
				break
			} else if role.roleType == "INVITED" {
				permissions.Type = "INVITED"
				break
			}
		}

		members = append(members, permissions)
	}

	_, err = tx.Exec(`
		DELETE FROM
			rooms_members
		WHERE
			room = $1;
	`, roomId)

	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
		INSERT INTO rooms_members (
			room,
			member,
			read,
			write,
			ban,
			invite,
			use_reply,
			use_secret,
			delete_other_message,
			type
		) VALUES (
			:room,
			:member,
			:read,
			:write,
			:ban,
			:invite,
			:use_reply,
			:use_secret,
			:delete_other_message,
			:type 
		)
	`, members)

	if err != nil {
		return err
	}

	return nil
}
