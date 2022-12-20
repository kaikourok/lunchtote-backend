package character

import (
	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) UpdateIcons(id int, icons *[]model.Icon, insertOnly bool) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		if !insertOnly {
			_, err := tx.Exec(`
				DELETE FROM
					characters_icons
				WHERE
					character = $1;
			`, id)
			if err != nil {
				return err
			}
		}

		if 0 < len(*icons) {
			type InsertCharacterIconStruct struct {
				Character int    `db:"character"`
				Path      string `db:"path"`
			}

			inserts := make([]InsertCharacterIconStruct, len(*icons))
			for i, icon := range *icons {
				inserts[i].Character = id
				inserts[i].Path = icon.Path
			}

			_, err := tx.NamedExec(`
				INSERT INTO characters_icons (
					character,
					path
				) VALUES (
					:character,
					:path
				)
			`, inserts)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
