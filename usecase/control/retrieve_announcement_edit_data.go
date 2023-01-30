package control

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *ControlUsecase) RetrieveAnnouncementEditData(announcementId int) (announcement *model.AnnouncementEditData, err error) {
	repository := s.registry.GetRepository()
	logger := s.registry.GetLogger()

	err = validation.Validate(announcementId, validation.Min(1))
	if err != nil {
		return nil, errors.ErrValidate
	}

	announcement, err = repository.RetrieveAnnouncementEditData(announcementId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return
}
