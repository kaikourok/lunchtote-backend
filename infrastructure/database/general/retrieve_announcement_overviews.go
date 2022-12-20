package general

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *GeneralRepository) RetrieveAnnouncementOverviews(basePoint, number int) (announcements *[]model.AnnouncementOverview, isContinue bool, err error) {
	rows, err := db.Queryx(`
		SELECT
			id,
			type,
			overview,
			announced_at,
			updated_at
		FROM
			announcements
		WHERE
			id < $1
		ORDER BY
			id DESC
		LIMIT
			$2;
	`, basePoint, number+1)
	if err != nil {
		return nil, false, err
	}

	announcementsSlice := make([]model.AnnouncementOverview, 0, number+1)
	for rows.Next() {
		var announcement model.AnnouncementOverview
		err := rows.Scan(
			&announcement.Id,
			&announcement.Type,
			&announcement.Overview,
			&announcement.AnnouncedAt,
			&announcement.UpdatedAt,
		)
		if err != nil {
			return nil, false, err
		}

		announcementsSlice = append(announcementsSlice, announcement)
	}

	isContinue = len(announcementsSlice) == number+1
	if isContinue {
		announcementsSlice = announcementsSlice[:len(announcementsSlice)-1]
	}

	return &announcementsSlice, isContinue, nil
}
