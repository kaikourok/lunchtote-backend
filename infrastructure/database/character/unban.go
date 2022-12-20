package character

import (
	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) Unban(targetId int, data *model.UnbanData) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`
			UPDATE
				characters
			SET
				banned_at = NULL
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
				'UNBAN'
			);
		`, targetId, data.Reason)
		if err != nil {
			return err
		}

		return nil
	})
}
