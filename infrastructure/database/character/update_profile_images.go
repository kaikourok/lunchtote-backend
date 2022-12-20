package character

import (
	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) UpdateProfileImages(id int, images *[]model.ProfileImage) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`
			DELETE FROM
				characters_icons
			WHERE
				character = $1;
		`, id)
		if err != nil {
			return err
		}

		if 0 < len(*images) {
			type InsertCharacterProfileImageStruct struct {
				Character int    `db:"character"`
				Path      string `db:"path"`
			}

			inserts := make([]InsertCharacterProfileImageStruct, len(*images))
			for i, profileImage := range *images {
				inserts[i].Character = id
				inserts[i].Path = profileImage.Path
			}

			_, err := tx.NamedExec(`
				INSERT INTO characters_profile_images (
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
