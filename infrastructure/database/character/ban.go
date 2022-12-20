package character

import (
	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) Ban(targetId int, data *model.BanData) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`
			UPDATE
				characters
			SET
				banned_at = CURRENT_TIMESTAMP
			WHERE
				id = $1;
		`, targetId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO character_prohibition_related_data (
				character,
				reason,
				type
			) VALUES (
				$1,
				$2,
				'BAN'
			);
		`, targetId, data.Reason)
		if err != nil {
			return err
		}

		return nil
	})
}
