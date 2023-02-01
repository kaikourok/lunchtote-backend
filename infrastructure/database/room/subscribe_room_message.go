package room

func (db *RoomRepository) SubscribeRoomMessage(characterId, roomId int) error {
	_, err := db.Exec(`
		INSERT INTO rooms_message_subscribers (
			character,
			room
		) VALUES (
			$1,
			$2
		);
	`, characterId, roomId)

	return err
}
