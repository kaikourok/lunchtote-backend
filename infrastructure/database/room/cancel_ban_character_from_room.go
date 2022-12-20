package room

import "errors"

func (db *RoomRepository) CancelBanCharacterFromRoom(userId, targetId, roomId int) error {
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

	var bannable, exists bool
	err := row.Scan(&bannable, &exists)

	if err != nil {
		return err
	}

	if !bannable {
		return errors.New("BAN権限がありません")
	}

	if !exists {
		return errors.New("指定のキャラクターはBANされていません")
	}

	_, err = db.Exec(`
		DELETE FROM
			rooms_banned_characters
		WHERE
			room   = $1 AND
			banned = $2;
	`, roomId, targetId)

	if err != nil {
		return err
	}

	return nil
}
