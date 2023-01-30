package general

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *GeneralRepository) RetrieveAnnouncementEditData(announcementId int) (announcement *model.AnnouncementEditData, err error) {
	row := db.QueryRowx(`
		SELECT
			type,
			title,
			overview,
			content,
			announced_at
		FROM
			announcements
		WHERE
			id = $1;
	`, announcementId)

	announcement = &model.AnnouncementEditData{}
	err = row.Scan(
		&announcement.Type,
		&announcement.Title,
		&announcement.Overview,
		&announcement.Content,
		&announcement.AnnouncedAt,
	)
	if err != nil {
		return nil, err
	}

	return
}
