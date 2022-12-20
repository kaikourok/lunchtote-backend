package room

import (
	"errors"
)

func (db *RoomRepository) DeleteRoom(characterId, roomId int) error {
	row := db.QueryRowx(`
		SELECT
			rooms.master
		FROM
			rooms
		WHERE
			rooms.id = $1;
	`, roomId)

	var masterId int
	err := row.Scan(&masterId)
	if err != nil {
		return err
	}

	if masterId != characterId {
		return errors.New("ルームを削除する権限がありません")
	}

	_, err = db.Exec(`
		UPDATE
			rooms
		SET
			deleted_at = CURRENT_TIMESTAMP
		WHERE
			id = $1;
	`, roomId)

	return err
}
