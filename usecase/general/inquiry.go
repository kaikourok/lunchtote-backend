package general

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *GeneralUsecase) Inquiry(characterId *int, inquiry string) error {
	repository := s.registry.GetRepository()
	logger := s.registry.GetLogger()

	err := validation.Validate(inquiry, validation.Required, validator.IsNotContainSpecialRune, validator.IsNotOnlySpace)
	if err != nil {
		return errors.ErrValidate
	}

	err = repository.Inquiry(characterId, inquiry)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
