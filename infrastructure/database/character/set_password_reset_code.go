package character

import (
	"github.com/jmoiron/sqlx"
)

func (db *CharacterRepository) SetPasswordResetCode(id int, email, code string, expireMinutes int) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		row := tx.QueryRowx(`
			SELECT
				true
			FROM
				characters
			WHERE
				id = $1 AND email = $2 AND deleted_at IS NULL;
		`, id, email)

		var exists bool
		err := row.Scan(&exists)
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

		_, err = tx.Exec(`
			INSERT INTO password_reset_codes (
				character,
				code,
				expire
			) VALUES (
				$1,
				$2,
				CURRENT_TIMESTAMP + ($3 * interval '1 minute')
			);
		`, id, code, expireMinutes)
		if err != nil {
			return err
		}

		return nil
	})
}
