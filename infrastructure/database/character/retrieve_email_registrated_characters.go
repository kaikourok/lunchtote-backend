package character

import (
	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) RetrieveEmailRegistratedCharacters(targetCharacters *[]int) (characters *[]model.CharacterEmailRegistratedData, err error) {
	sql := `
		SELECT
			id,
			email
		FROM
			characters
		WHERE
			email IS NOT NULL AND deleted_at IS NULL AND banned_at IS NULL
	`

	params := make([]interface{}, 0)
	if targetCharacters == nil {
		sql += `;`
	} else {
		sql += `
			AND id IN (?);
		`
		sql, params, err = sqlx.In(sql, *targetCharacters)
		if err != nil {
			return nil, err
		}
	}

	rows, err := db.Queryx(sql, params...)
	if err != nil {
		return nil, err
	}

	charactersSlice := make([]model.CharacterEmailRegistratedData, 0, 2048)
	for rows.Next() {
		var character model.CharacterEmailRegistratedData
		err = rows.Scan(&character.Id, &character.Email)
		if err != nil {
			return nil, err
		}

		charactersSlice = append(charactersSlice, character)
	}

	return &charactersSlice, nil
}
