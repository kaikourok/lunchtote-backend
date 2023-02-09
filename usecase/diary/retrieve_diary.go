package diary

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/entity/service"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *DiaryUsecase) RetrieveDiary(characterId *int, targetId, nth int) (*model.Diary, error) {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	err := validation.Validate(targetId, validation.Min(1))
	if err != nil {
		return nil, errors.ErrValidate
	}

	err = validation.Validate(nth, validation.Min(1))
	if err != nil {
		return nil, errors.ErrValidate
	}

	diary, err := repository.RetrieveDiary(characterId, targetId, nth)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	diary.Diary = service.StylizeTextEntry(diary.Diary, config.GetString("general.upload-directory"))

	return diary, nil
}
