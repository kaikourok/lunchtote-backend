package room

func (db *RoomRepository) SubscribeRoomNewMember(characterId, roomId int) error {
	_, err := db.Exec(`
		INSERT INTO rooms_new_member_subscribers (
			character,
			room
		) VALUES (
			$1,
			$2
		);
	`, characterId, roomId)

	return err
}
