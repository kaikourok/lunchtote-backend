package character

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

func (db *CharacterRepository) Mute(userId, targetId int) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		row := tx.QueryRowx(`
			SELECT
				COALESCE((SELECT true FROM mutes      WHERE muter = $1 AND muted = $2), false),
				COALESCE((SELECT true FROM characters WHERE deleted_at IS NULL AND administrator = false AND id = $2), false);
		`, userId, targetId)

		var mute, targetExists bool
		err := row.Scan(&mute, &targetExists)
		if err != nil {
			return err
		}

		if !targetExists {
			err := errors.New("対象が存在しません")
			return err
		}

		if mute {
			err := errors.New("すでにミュートしています")
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO	mutes (muter, muted) VALUES ($1, $2);
		`, userId, targetId)
		if err != nil {
			return err
		}

		return nil
	})
}
