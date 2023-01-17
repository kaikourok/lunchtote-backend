package room

func (db *RoomRepository) RenameRoomMessageFetchConfig(characterId, configId int, name string) error {
	_, err := db.Exec(`
		UPDATE
			message_fetch_configs
		SET
			name = $3
		WHERE
			id     = $1 AND
			master = $2;
	`, configId, characterId, name)

	return err
}
