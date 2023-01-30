package general

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *GeneralRepository) UpdateAnnouncement(announcementId int, announce *model.AnnouncementEditDataUpdate) error {
	_, err := db.Exec(`
		UPDATE
			announcements
		SET
			type         = $2,
			title        = $3,
			overview     = $4,
			content      = $5,
			updated_at = (
				CASE
					WHEN $6 = true THEN
						updated_at
					ELSE
						CURRENT_TIMESTAMP
				END
			)
		WHERE
			id = $1;
	`,
		announcementId,
		announce.Type,
		announce.Title,
		announce.Overview,
		announce.Content,
		announce.SilentUpdate,
	)

	return err
}
