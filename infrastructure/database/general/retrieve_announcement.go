package general

import (
	"database/sql"

	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *GeneralRepository) RetrieveAnnouncement(announcementId int) (announcement *model.Announcement, prevGuide, nextGuide *model.AnnouncementGuideData, err error) {
	retrieveGuideData := func(isPrev bool) (*model.AnnouncementGuideData, error) {
		statement := `
			SELECT
				id,
				title
			FROM
				announcements
		`

		if isPrev {
			statement += `
				WHERE
					id < $1
				ORDER BY
					id DESC
			`
		} else {
			statement += `
				WHERE
					$1 < id
				ORDER BY
					id ASC
			`
		}

		statement += ` 
			LIMIT
				1;
		`

		row := db.QueryRowx(statement, announcementId)

		var data model.AnnouncementGuideData
		err := row.Scan(&data.Id, &data.Title)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			} else {
				return nil, err
			}
		}

		return &data, nil
	}

	row := db.QueryRowx(`
		SELECT
			id,
			type,
			title,
			content,
			announced_at,
			updated_at
		FROM
			announcements
		WHERE
			id = $1;
	`, announcementId)

	announcement = &model.Announcement{}
	err = row.Scan(
		&announcement.Id,
		&announcement.Type,
		&announcement.Title,
		&announcement.Content,
		&announcement.AnnouncedAt,
		&announcement.UpdatedAt,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	prevGuide, err = retrieveGuideData(true)
	if err != nil {
		return nil, nil, nil, err
	}

	nextGuide, err = retrieveGuideData(false)
	if err != nil {
		return nil, nil, nil, err
	}

	return announcement, prevGuide, nextGuide, nil
}
