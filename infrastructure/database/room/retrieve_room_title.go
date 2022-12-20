package room

func (db *RoomRepository) RetrieveRoomTitle(roomId int) (title string, err error) {
	row := db.QueryRowx(`
		SELECT
			title
		FROM
			rooms
		WHERE
			id = $1;
	`, roomId)

	err = row.Scan(&title)
	return
}
