package character

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

func (db *CharacterRepository) Block(userId, targetId int) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		row := tx.QueryRowx(`
			SELECT
				COALESCE((SELECT true FROM blocks     WHERE blocker = $1 AND blocked = $2), false),
				COALESCE((SELECT true FROM characters WHERE deleted_at IS NULL AND administrator = false AND id = $2), false);
		`, userId, targetId)

		var block, targetExists bool
		err := row.Scan(&block, &targetExists)
		if err != nil {
			return err
		}

		if !targetExists {
			err := errors.New("対象が存在しません")
			return err
		}

		if block {
			err := errors.New("すでにブロックしています")
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO	blocks
				(blocker, blocked)
			VALUES
				($1, $2);
		`, userId, targetId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			DELETE FROM
				follows
			WHERE
				(follower = $1 AND followed = $2) OR
				(follower = $2 AND followed = $1);
		`, userId, targetId)
		if err != nil {
			return err
		}

		return nil
	})
}
