package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) RegisterGoogleData(characterId int, googleId string) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(googleId, validation.Required)
	if err != nil {
		return errors.ErrValidate
	}

	err = repository.RegisterGoogleData(characterId, googleId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
