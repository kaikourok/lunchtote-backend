package general

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *GeneralRepository) Announce(announce *model.AnnouncementEditData) error {
	_, err := db.Exec(`
		INSERT INTO announcements (
			type,
			title,
			overview,
			content
		) VALUES (
			$1,
			$2,
			$3,
			$4
		);
	`,
		announce.Type,
		announce.Title,
		announce.Overview,
		announce.Content,
	)

	return err
}
