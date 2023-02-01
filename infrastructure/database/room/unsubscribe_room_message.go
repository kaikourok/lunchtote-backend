package room

func (db *RoomRepository) UnsubscribeRoomMessage(characterId, roomId int) error {
	_, err := db.Exec(`
		DELETE FROM
			rooms_message_subscribers
		WHERE
			character = $1 AND
			room      = $2;
	`, characterId, roomId)

	return err
}
