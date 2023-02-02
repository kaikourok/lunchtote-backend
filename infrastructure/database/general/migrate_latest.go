package general

func (db *GeneralRepository) MigrateLatest() error {
	migrate, err := db.createMigrateInstance()
	if err != nil {
		return err
	}

	return migrate.Up()
}
