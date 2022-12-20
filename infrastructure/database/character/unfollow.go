package character

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

func (db *CharacterRepository) Unfollow(userId, targetId int) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		row := tx.QueryRowx(`
			SELECT
				COALESCE((SELECT true FROM follows    WHERE follower = $1 AND followed = $2), false),
				COALESCE((SELECT true FROM characters WHERE deleted_at IS NULL AND administrator = false AND id = $2), false);
		`, userId, targetId)

		var follow, targetExists bool
		err := row.Scan(&follow, &targetExists)
		if err != nil {
			return err
		}

		if !targetExists {
			err := errors.New("対象が存在しません")
			return err
		}

		if !follow {
			err := errors.New("フォローしていません")
			return err
		}

		_, err = tx.Exec(`
			DELETE FROM follows WHERE follower = $1 AND followed = $2;
		`, userId, targetId)
		if err != nil {
			return err
		}

		return nil
	})
}
