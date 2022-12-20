package diary

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *DiaryUsecase) ReservePublishDiary(characterId int, title string, diary string) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(title, validation.Required, validator.IsNotContainSpecialRune, validator.IsNotOnlySpace)
	if err != nil {
		return errors.ErrValidate
	}

	err = validation.Validate(diary, validation.Required, validator.IsNotContainSpecialRune, validator.IsNotOnlySpace)
	if err != nil {
		return errors.ErrValidate
	}

	err = repository.ReservePublishDiary(characterId, title, diary)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
