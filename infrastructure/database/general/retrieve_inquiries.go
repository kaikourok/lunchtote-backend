package general

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *GeneralRepository) RetrieveInquiries(basePoint, number int, unresolvedOnly bool) (inquiries *[]model.Inquiry, isContinue bool, err error) {
	sql := `
		SELECT
			inquiries.id,
			inquiries.inquiry,
			inquiries.posted_at,
			inquiries.resolved_at IS NOT NULL,
			characters.id,
			characters.name,
			characters.mainicon
		FROM
			inquiries
		LEFT JOIN
			characters ON (inquiries.character = characters.id)
		WHERE
			inquiries.id < $1
	`

	if unresolvedOnly {
		sql += `
			AND resolved_at IS NULL
		`
	}

	sql += `
		ORDER BY
			inquiries.id DESC
		LIMIT
			$2;
	`

	rows, err := db.Queryx(sql, basePoint, number+1)
	if err != nil {
		return nil, false, err
	}

	inquiriesSlice := make([]model.Inquiry, 0, number+1)
	var characterId *int
	var characterName, characterMainicon *string
	for rows.Next() {
		var inquiry model.Inquiry
		err := rows.Scan(
			&inquiry.Id,
			&inquiry.Content,
			&inquiry.PostedAt,
			&inquiry.Resolved,
			&characterId,
			&characterName,
			&characterMainicon,
		)
		if err != nil {
			return nil, false, err
		}

		if characterId != nil && characterName != nil && characterMainicon != nil {
			inquiry.Character = &model.CharacterOverview{
				Id:       *characterId,
				Name:     *characterName,
				Mainicon: *characterMainicon,
			}
		}

		inquiriesSlice = append(inquiriesSlice, inquiry)
	}

	isContinue = len(inquiriesSlice) == number+1
	if isContinue {
		inquiriesSlice = inquiriesSlice[:len(inquiriesSlice)-1]
	}

	return &inquiriesSlice, isContinue, nil
}
