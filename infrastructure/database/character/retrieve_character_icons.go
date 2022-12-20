package character

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *CharacterRepository) RetrieveCharacterIcons(id int) (icons *[]model.Icon, err error) {
	rows, err := db.Queryx(`
		SELECT
			path
		FROM
			characters_icons
		WHERE
			character = $1;
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	iconsSlice := make([]model.Icon, 0, 64)
	for rows.Next() {
		var icon model.Icon
		err = rows.Scan(&icon.Path)
		if err != nil {
			return nil, err
		}

		iconsSlice = append(iconsSlice, icon)
	}

	return &iconsSlice, nil
}
