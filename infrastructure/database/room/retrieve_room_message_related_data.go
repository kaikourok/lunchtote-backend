package room

func (db *RoomRepository) RetrieveRoomMessageRelatedData(messageId int) (roomId, senderId int, err error) {
	row := db.QueryRowx(`
		SELECT
			room,
			character
		FROM
			rooms_messages
		WHERE
			id = $1;
	`, messageId)

	err = row.Scan(&roomId, &senderId)
	return
}
