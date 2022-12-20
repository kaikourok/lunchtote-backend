package character

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

func (db *CharacterRepository) UpdatePasswordByResetCode(id int, code, password string) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		row := tx.QueryRowx(`
			SELECT
				COALESCE((
					SELECT
						true
					FROM
						password_reset_codes
					JOIN
						characters ON (password_reset_codes.character = characters.id)
					WHERE
						password_reset_codes.character = $1 AND
						password_reset_codes.code      = $2 AND
						CURRENT_TIMESTAMP < password_reset_codes.expire AND
						characters.deleted_at IS NULL
				), false)
		`, id, code)

		var exists bool
		err := row.Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			err := errors.New("NOT_FOUND")
			return err
		}

		_, err = tx.Exec(`
			UPDATE
				characters
			SET
				password = $2
			WHERE
				id = $1;
		`, id, password)

		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			DELETE FROM
				password_reset_codes
			WHERE
				character = $1;
		`, id)
		if err != nil {
			return err
		}

		return nil
	})
}
