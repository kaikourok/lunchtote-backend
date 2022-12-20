package character

import (
	"github.com/jmoiron/sqlx"
)

func (db *CharacterRepository) CreateAdministrator(id int, hashedPassword, name, nickname, username, notificationToken string) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`
			INSERT INTO characters (
				id,
				password,
				name,
				nickname,
				username,
				administrator,
				notification_token
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7
			);
		`,
			id,
			hashedPassword,
			name,
			nickname,
			username,
			true,
			notificationToken,
		)
		if err != nil {
			return err
		}

		return nil
	})
}
