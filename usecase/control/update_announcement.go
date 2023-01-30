package control

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *ControlUsecase) UpdateAnnouncement(announcementId int, announce *model.AnnouncementEditDataUpdate) error {
	repository := s.registry.GetRepository()
	logger := s.registry.GetLogger()

	err := validation.Validate(announcementId, validation.Min(1))
	if err != nil {
		return errors.ErrValidate
	}

	err = repository.UpdateAnnouncement(announcementId, announce)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
