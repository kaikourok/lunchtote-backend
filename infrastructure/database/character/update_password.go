package character

func (db *CharacterRepository) UpdatePassword(id int, password string) error {
	_, err := db.Exec(`
		UPDATE
			characters
		SET
			password = $2
		WHERE
			id = $1 AND deleted_at IS NULL;
	`, id, password)

	return err
}
