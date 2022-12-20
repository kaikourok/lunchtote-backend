package database

import (
	"html"

	"github.com/jmoiron/sqlx"
)

type SqlxDatabase interface {
	Beginx() (*sqlx.Tx, error)
}

func ExecTx(db SqlxDatabase, f func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	html.EscapeString("aaa")

	err = f(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
