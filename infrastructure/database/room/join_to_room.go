package room

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

func (db *RoomRepository) JoinToRoom(targetId, roomId int) error {
	row := db.QueryRowx(`
		SELECT
			EXISTS (
				SELECT
					*
				FROM
					rooms_members
				WHERE
					room   = $1 AND
					member = $2 
			)
		;
	`, roomId, targetId)

	var isAlreadyExists bool
	err := row.Scan(&isAlreadyExists)
	if err != nil {
		return err
	}

	if isAlreadyExists {
		return errors.New("すでにメンバーになっているキャラクターです")
	}

	return db.ExecTx(func(tx *sqlx.Tx) error {
		_, err = tx.Exec(`
			INSERT INTO rooms_members (
				room,
				member,
				write,
				ban,
				invite,
				use_reply,
				use_secret,
				delete_other_message,
				create_children_room,
				type
			)

			SELECT
				$1,
				$2,
				write,
				ban,
				invite,
				use_reply,
				use_secret,
				delete_other_message,
				create_children_room,
				'MEMBER'
			FROM
				rooms_roles
			WHERE
				rooms_roles.room = $1 AND
				rooms_roles.type = 'DEFAULT';
		`, roomId, targetId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO rooms_roles_members (
				role,
				character
			)

			SELECT
				id,
				$2
			FROM
				rooms_roles
			WHERE
				rooms_roles.room = $1 AND
				rooms_roles.type = 'DEFAULT';
		`, roomId, targetId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE
				rooms
			SET
				members_count = members_count + 1
			WHERE
				id = $1;
		`, roomId)
		if err != nil {
			return err
		}

		return nil
	})

}
