package character

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

func (db *CharacterRepository) DeleteList(userId, listId int) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		row := tx.QueryRowx(`
			SELECT
				master
			FROM
				lists
			WHERE
				id = $1;
		`, listId)

		var listMaster int
		err := row.Scan(&listMaster)
		if err != nil {
			return err
		}

		if userId != listMaster {
			err := errors.New("リストの管理者ではありません")
			return err
		}

		_, err = tx.Exec(`
			DELETE FROM
				lists_characters
			WHERE
				list = $1;
		`, listId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			DELETE FROM
				lists
			WHERE
				id = $1;
		`, listId)

		return err
	})
}
