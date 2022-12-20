package character

import (
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
)

// TODO: 明らかに肥大化しすぎなのでdomainとinterface/storageに適切に分割する
func (db *CharacterRepository) DeleteUploadedImages(characterId int, imageIds *[]int, uploadDir string) error {
	query := `
		DELETE FROM
			characters_uploaded_images
		WHERE
			character = ` + strconv.Itoa(characterId) + ` AND
			id IN (?)
		RETURNING
			path;
	`

	query, args, err := sqlx.In(query, *imageIds)
	if err != nil {
		return err
	}

	rows, err := db.Queryx(db.Rebind(query), args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var path string
		err = rows.Scan(&path)
		if err != nil {
			return err
		}

		os.Remove(uploadDir + "/" + path)
	}

	return nil
}
