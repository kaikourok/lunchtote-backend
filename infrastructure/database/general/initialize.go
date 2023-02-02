package general

func (db *GeneralRepository) Initialize() error {
	db.MustExec(`
		INSERT INTO game_status (
			nth
		) VALUES (
			0
		);
  `)

	return nil
}
