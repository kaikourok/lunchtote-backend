package general

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *GeneralUsecase) RetrieveAnnouncement(announcementId int) (announcement *model.Announcement, prevGuide, nextGuide *model.AnnouncementGuideData, err error) {
	repository := s.registry.GetRepository()
	logger := s.registry.GetLogger()

	err = validation.Validate(announcementId, validation.Min(1))
	if err != nil {
		return nil, nil, nil, errors.ErrValidate
	}

	announcement, prevGuide, nextGuide, err = repository.RetrieveAnnouncement(announcementId)
	if err != nil {
		logger.Error(err)
		return nil, nil, nil, err
	}

	return
}
