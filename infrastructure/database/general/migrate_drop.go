package general

func (db *GeneralRepository) MigrateDrop() error {
	_, err := db.Exec(`
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
	`)
	return err
}
