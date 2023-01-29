package character

import (
	"github.com/jmoiron/sqlx"
)

func (db *CharacterRepository) CreateCharacter(name, nickname, username, password, notificationToken string) (id int, err error) {
	err = db.ExecTx(func(tx *sqlx.Tx) error {
		row := tx.QueryRowx(`
			INSERT INTO characters (
				name,
				nickname,
				username,
				password,
				notification_token
			) VALUES (
				$1, $2, $3, $4, $5
			)
			RETURNING
				seq;
		`,
			name,
			nickname,
			username,
			password,
			notificationToken,
		)

		var seq int
		err := row.Scan(&seq)
		if err != nil {
			return err
		}

		row = tx.QueryRowx(`
			UPDATE characters SET
				id = (SELECT count(*) FROM characters WHERE administrator = false AND seq <= $1)
			WHERE
				seq = $1
			RETURNING
				id;
		`, seq)

		err = row.Scan(&id)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO message_fetch_configs (
				master,
				config_order,
				name,
				category,
				relate_filter,
				children
			) VALUES (
				$1,
				0,
				'全体',
				'all',
				true,
				true
			), (
				$1,
				1,
				'フォロー',
				'follow',
				true,
				true
			), (
				$1,
				2,
				'返信',
				'replied',
				true,
				true
			);
		`, id)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
