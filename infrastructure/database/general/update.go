package general

import "github.com/jmoiron/sqlx"

func (db *GeneralRepository) Update() error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`
			INSERT INTO diaries (
				character,
				title,
				diary,
				nth
			)
			
			SELECT
				id           AS character,
				diary_title  AS title,
				diary        AS diary,
				(SELECT nth FROM game_status) AS nth
			FROM
				characters
			WHERE
				diary_title IS NOT NULL AND
				diary       IS NOT NULL;
		`)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE
				characters
			SET
				diary_title = NULL,
				diary       = NULL;
		`)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE
				game_status
			SET
				nth = nth + 1;
		`)

		return err
	})
}
