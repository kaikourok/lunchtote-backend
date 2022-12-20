package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/lib/pq"
)

func (db *CharacterRepository) RetrieveProfileEditData(id int) (data *model.ProfileEditData, err error) {
	row := db.QueryRowx(`
		SELECT
			name,
			nickname,
			summary,
			profile,
			Mainicon,
			ARRAY(SELECT tag FROM characters_tags WHERE character = $1)
		FROM
			characters
		WHERE
			deleted_at IS NULL AND id = $1;
	`, id)

	data = &model.ProfileEditData{}
	err = row.Scan(
		&data.Name,
		&data.Nickname,
		&data.Summary,
		&data.Profile,
		&data.Mainicon,
		pq.Array(&data.Tags),
	)
	return
}
