package room

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

func (db *RoomRepository) BanCharacterFromRoom(userId, targetId, roomId int) error {
	row := db.QueryRowx(`
		SELECT
			EXISTS (
				SELECT
					*
				FROM
					rooms_members
				WHERE
					room   = $1 AND
					member = $2 AND
					ban    = true
			),
			EXISTS (
				SELECT
					*
				FROM
					rooms_members
				WHERE
					room   = $1 AND
					member = $3 AND
					type != 'MASTER'
			);
	`, roomId, userId, targetId)

	var bannable, existsTargetMember bool
	err := row.Scan(&bannable, &existsTargetMember)

	if err != nil {
		return err
	}

	if !bannable {
		return errors.New("BANを行う権限がありません")
	}

	if !existsTargetMember {
		return errors.New("対象はメンバーではありません")
	}

	return db.ExecTx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`
			DELETE FROM
				rooms_members
			WHERE
				room   = $1 AND
				member = $2;
		`, roomId, targetId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			DELETE FROM
				rooms_roles_members
			WHERE
				role IN (
					SELECT
						id
					FROM
						rooms_roles
					WHERE
						room = $1
				) AND
				character = $2;
		`, roomId, targetId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO rooms_banned_characters (
				room,
				banned,
				banner
			) VALUES (
				$1,
				$2,
				$3
			);
		`, roomId, targetId, userId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE
				rooms
			SET
				members_count = members_count - 1
			WHERE
				id = $1;
		`, roomId)
		if err != nil {
			return err
		}

		return nil
	})
}
