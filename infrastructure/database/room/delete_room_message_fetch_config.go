package room

func (db *RoomRepository) DeleteRoomMessageFetchConfig(characterId, configId int) error {
	_, err := db.Exec(`
		DELETE FROM
			message_fetch_configs
		WHERE
			id     = $1 AND
			master = $2;
	`, configId, characterId)

	return err
}
