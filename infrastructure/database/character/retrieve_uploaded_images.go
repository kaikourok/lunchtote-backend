package character

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *CharacterRepository) RetrieveUploadedImages(id int) (images *[]model.UploadedImage, err error) {
	rows, err := db.Queryx(`
		SELECT
			id,
			path,
			uploaded_at
		FROM
			characters_uploaded_images
		WHERE
			character = $1;
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	uploadedImages := make([]model.UploadedImage, 0, 64)
	for rows.Next() {
		var uploadedImage model.UploadedImage
		err = rows.Scan(
			&uploadedImage.Id,
			&uploadedImage.Path,
			&uploadedImage.UploadedAt,
		)
		if err != nil {
			return nil, err
		}

		uploadedImages = append(uploadedImages, uploadedImage)
	}

	return &uploadedImages, nil
}
