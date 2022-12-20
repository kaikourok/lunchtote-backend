package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) RetrieveProhibitionRelatedData(targetId int) (data *[]model.ProhibitionRelatedData, err error) {
	rows, err := db.Queryx(`
		SELECT
			id,
			type,
			reason,
			timestamp
		FROM 
			character_prohibition_related_data
		WHERE
			character = $1
		ORDER BY
			timestamp DESC;
	`, targetId)
	if err != nil {
		return nil, err
	}

	relateds := make([]model.ProhibitionRelatedData, 0, 8)
	for rows.Next() {
		var related model.ProhibitionRelatedData
		err = rows.Scan(
			&related.Id,
			&related.Type,
			&related.Reason,
			&related.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		relateds = append(relateds, related)
	}

	return &relateds, nil
}
