package room

import "errors"

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
					rooms_banned_characters
				WHERE
					room   = $1 AND
					banned = $3 
			);
	`, roomId, userId, targetId)

	var bannable, isAlreadyExists bool
	err := row.Scan(&bannable, &isAlreadyExists)

	if err != nil {
		return err
	}

	if !bannable {
		return errors.New("BANを行う権限がありません")
	}

	if isAlreadyExists {
		return errors.New("すでにBANしているキャラクターです")
	}

	_, err = db.Exec(`
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

	return err
}
