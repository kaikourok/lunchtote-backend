package character

import (
	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) Adomonish(targetId int, data *model.AdomonishData) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`
			INSERT INTO mails (
				receiver,
				name,
				title,
				message
			) VALUES (
				$1,
				'管理者',
				$2,
				$3
			);
		`,
			targetId,
			data.Title,
			data.Message,
		)
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
				'ADOMONISH'
			);
		`, targetId, data.Reason)
		if err != nil {
			return err
		}

		return nil
	})
}
