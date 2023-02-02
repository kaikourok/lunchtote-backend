package general

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (db *GeneralRepository) createMigrateInstance() (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db.DB.DB, &postgres.Config{})
	if err != nil {

	}

	return migrate.NewWithDatabaseInstance(
		"file://infrastructure/database/migrations",
		"postgres", driver)
}
