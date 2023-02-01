package character

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

func (db *CharacterRepository) Follow(userId, targetId int) (userName string, webhook string, err error) {
	err = db.ExecTx(func(tx *sqlx.Tx) error {
		row := tx.QueryRowx(`
			SELECT
				COALESCE((SELECT true FROM follows    WHERE follower = $1 AND followed = $2), false),
				COALESCE((SELECT true FROM blocks     WHERE blocker  = $1 AND blocked  = $2), false),
				COALESCE((SELECT true FROM blocks     WHERE blocked  = $1 AND blocker  = $2), false),
				COALESCE((SELECT true FROM characters WHERE deleted_at IS NULL AND administrator = false AND id = $2), false),
				(SELECT notification_followed FROM characters WHERE id = $2),
				(SELECT nickname              FROM characters WHERE id = $1),
				COALESCE((SELECT webhook FROM characters WHERE id = $2 AND webhook_followed = true), '');
		`, userId, targetId)

		var follow, block, blocked, targetExists, notificationFollowed bool
		var userNickname, targetWebhook string
		err := row.Scan(
			&follow,
			&block,
			&blocked,
			&notificationFollowed,
			&targetExists,
			&userNickname,
			&targetWebhook,
		)
		if err != nil {
			return err
		}

		if !targetExists {
			err := errors.New("対象が存在しません")
			return err
		}

		if follow {
			err := errors.New("すでにフォローしています")
			return err
		}

		if block || blocked {
			err := errors.New("ブロックしている、あるいはブロックされています")
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO	follows (follower, followed) VALUES ($1, $2);
		`, userId, targetId)
		if err != nil {
			return err
		}

		if notificationFollowed {
			row = tx.QueryRowx(`
				INSERT INTO notifications (
					character,
					type
				) VALUES (
					$1,
					'FOLLOWED'
				)

				RETURNING
					id;
			`, targetId)
			var notificationId int
			err = row.Scan(&notificationId)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`
				INSERT INTO notifications_followed_data (
					notification,
					followed_by
				) VALUES (
					$1,
					$2
				);
			`, notificationId, userId)
			if err != nil {
				return err
			}
		}

		userName = userNickname
		webhook = targetWebhook
		return nil
	})

	return
}
