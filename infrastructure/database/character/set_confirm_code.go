package character

import (
	"github.com/jmoiron/sqlx"
)

func (db *CharacterRepository) SetConfirmCode(id int, email, code string, expireMinutes int) error {
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
			INSERT INTO mail_confirm_codes (
				character,
				email,
				code,
				expire
			) VALUES (
				$1,
				$2,
				$3,
				CURRENT_TIMESTAMP + ($4 * interval '1 minute')
			);
		`, id, email, code, expireMinutes)

		return err
	})
}
