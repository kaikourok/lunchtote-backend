package character

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *CharacterRepository) RetrieveLists(id int) (lists *[]model.ListOverview, err error) {
	rows, err := db.Queryx(`
		SELECT
			lists.id,
			lists.name
		FROM
			lists
		WHERE
			lists.master = $1;
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	listOverviews := make([]model.ListOverview, 0, 64)
	for rows.Next() {
		var list model.ListOverview
		err = rows.Scan(
			&list.Id,
			&list.Name,
		)
		if err != nil {
			return nil, err
		}

		listOverviews = append(listOverviews, list)
	}

	return &listOverviews, nil
}
