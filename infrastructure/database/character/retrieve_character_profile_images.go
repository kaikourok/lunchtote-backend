package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) RetrieveCharacterProfileImages(id int) (images *[]model.ProfileImage, err error) {
	rows, err := db.Queryx(`
		SELECT
			path
		FROM
			characters_profile_images
		WHERE
			character = $1;
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profileImages := make([]model.ProfileImage, 0, 64)
	for rows.Next() {
		var profileImage model.ProfileImage
		err = rows.Scan(&profileImage.Path)
		if err != nil {
			return nil, err
		}

		profileImages = append(profileImages, profileImage)
	}

	return &profileImages, nil
}
