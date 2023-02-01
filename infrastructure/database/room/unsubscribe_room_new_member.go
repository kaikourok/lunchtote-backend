package room

func (db *RoomRepository) UnsubscribeRoomNewMember(characterId, roomId int) error {
	_, err := db.Exec(`
		DELETE FROM
			rooms_new_member_subscribers
		WHERE
			character = $1 AND
			room      = $2;
	`, characterId, roomId)

	return err
}
