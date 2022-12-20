package character

import (
	"github.com/jmoiron/sqlx"
)

func (db *CharacterRepository) UpdateEmail(id int, email string) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`
			DELETE FROM
				mail_confirm_codes
			WHERE
				character = $1;
		`, id)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE
				characters
			SET
				email = $2
			WHERE
				id = $1 AND deleted_at IS NULL;
		`, id, email)

		return err
	})
}
