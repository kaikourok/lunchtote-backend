package room

func (db *RoomRepository) DeleteRoomMessage(messageId int) error {
	_, err := db.Exec(`
		UPDATE
			rooms_messages
		SET
			deleted_at = CURRENT_TIMESTAMP
		WHERE
			id = $1;
	`, messageId)

	return err
}
