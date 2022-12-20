package room

import "errors"

func (db *RoomRepository) InviteCharacterToRoom(userId, targetId, roomId int) error {
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
					invite = true
			),
			EXISTS (
				SELECT
					*
				FROM
					rooms_members
				WHERE
					room   = $1 AND
					member = $3 
			);
	`, roomId, userId, targetId)

	var invitable, isAlreadyExists bool
	err := row.Scan(&invitable, &isAlreadyExists)

	if err != nil {
		return err
	}

	if !invitable {
		return errors.New("招待を行う権限がありません")
	}

	if isAlreadyExists {
		return errors.New("すでにメンバーになっているキャラクターです")
	}

	_, err = db.Exec(`
		INSERT INTO rooms_invited_characters (
			room,
			invited,
			inviter
		) VALUES (
			$1,
			$2,
			$3
		);
	`, roomId, targetId, userId)

	if err != nil {
		return err
	}

	return nil
}
