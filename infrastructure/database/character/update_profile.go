package character

import (
	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) UpdateProfile(id int, profile *model.ProfileEditData) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`
			UPDATE
				characters
			SET
				name     = $2,
				nickname = $3,
				summary  = $4,
				profile  = $5,
				mainicon = $6
			WHERE
				id = $1;
		`,
			id,
			profile.Name,
			profile.Nickname,
			profile.Summary,
			profile.Profile,
			profile.Mainicon,
		)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			DELETE FROM
				characters_tags
			WHERE
				character = $1;
		`, id)
		if err != nil {
			return err
		}

		if 0 < len(profile.Tags) {
			type InsertCharacterTagStruct struct {
				Character int    `db:"character"`
				Tag       string `db:"tag"`
			}

			inserts := make([]InsertCharacterTagStruct, len(profile.Tags))
			for i, tag := range profile.Tags {
				inserts[i] = InsertCharacterTagStruct{
					Character: id,
					Tag:       tag,
				}
			}

			_, err = db.NamedExec(`
				INSERT INTO characters_tags (
					character,
					tag
				) VALUES (
					:character,
					:tag
				)
			`, inserts)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
