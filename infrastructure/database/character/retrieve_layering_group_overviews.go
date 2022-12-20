package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) RetrieveLayeringGroupOverviews(characterId int) (overviews *[]model.CharacterIconLayeringGroupOverview, err error) {
	rows, err := db.Queryx(`
		SELECT
			id,
			name
		FROM
			characters_icon_layering_groups
		WHERE
			character = $1;
	`, characterId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	overviewsSlice := make([]model.CharacterIconLayeringGroupOverview, 0, 16)
	for rows.Next() {
		var overview model.CharacterIconLayeringGroupOverview
		err = rows.Scan(&overview.Id, &overview.Name)
		if err != nil {
			return nil, err
		}

		overviewsSlice = append(overviewsSlice, overview)
	}

	return &overviewsSlice, nil
}
