package mail

import (
	"github.com/jmoiron/sqlx"
)

func (db *MailRepository) DeleteMails(characterId int, mailIds *[]int) (deletedIds *[]int, err error) {
	deleteds := make([]int, 0, len(*mailIds))

	err = db.ExecTx(func(tx *sqlx.Tx) error {
		query := `
			UPDATE
				mails
			SET
				deleted_at = CURRENT_TIMESTAMP
			WHERE
				id IN (:idList) AND receiver = :receiver AND deleted_at IS NULL
			RETURNING
				id;
		`

		query, args, err := sqlx.Named(query, map[string]interface{}{
			"idList":   *mailIds,
			"receiver": characterId,
		})
		if err != nil {
			return err
		}

		query, args, err = sqlx.In(query, args...)
		if err != nil {
			return err
		}
		query = db.Rebind(query)

		rows, err := tx.Queryx(query, args...)
		if err != nil {
			return err
		}

		for rows.Next() {
			var deletedId int
			err = rows.Scan(&deletedId)
			if err != nil {
				return err
			}

			deleteds = append(deleteds, deletedId)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &deleteds, err
}
