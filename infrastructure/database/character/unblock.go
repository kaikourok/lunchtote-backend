package character

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

func (db *CharacterRepository) Unblock(userId, targetId int) error {
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

		if !block {
			err := errors.New("ブロックしていません")
			return err
		}

		_, err = tx.Exec(`
			DELETE FROM blocks WHERE blocker = $1 AND blocked = $2;
		`, userId, targetId)
		if err != nil {
			return err
		}

		return nil
	})
}
