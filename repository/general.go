package repository

import "github.com/kaikourok/lunchtote-backend/entity/model"

type generalRepository interface {
	Initialize() error
	Update() error
	Connections() (connections int, err error)
	Inquiry(characterId *int, inquiry string) error
	Announce(announce *model.AnnouncementData) error
	UpdateAnnouncement(announceId int, announce *model.AnnouncementEditDataUpdate) error
	UpdateInquiryState(inquiryId int, resolved bool) error
	RetrieveInquiries(basePoint, number int, unresolvedOnly bool) (inquiries *[]model.Inquiry, isContinue bool, err error)
	RetrieveAnnouncement(announcementId int) (announcement *model.Announcement, prevGuide, nextGuide *model.AnnouncementGuideData, err error)
	RetrieveAnnouncementEditData(announcementId int) (announcement *model.AnnouncementEditData, err error)
	RetrieveAnnouncementOverviews(basePoint, number int) (announcements *[]model.AnnouncementOverview, isContinue bool, err error)
}
