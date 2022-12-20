package room

import "errors"

func (db *RoomRepository) CancelInviteCharacterToRoom(userId, targetId, roomId int) error {
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
					rooms_invited_characters 
				WHERE
					room    = $1 AND
					invited = $3 
			);
	`, roomId, userId, targetId)

	var invitable, exists bool
	err := row.Scan(&invitable, &exists)

	if err != nil {
		return err
	}

	if !invitable {
		return errors.New("招待権限がありません")
	}

	if !exists {
		return errors.New("指定のキャラクターには招待を行っていません")
	}

	_, err = db.Exec(`
		DELETE FROM
			rooms_invited_characters
		WHERE
			room    = $1 AND
			invited = $2;
	`, roomId, targetId)

	return err
}
