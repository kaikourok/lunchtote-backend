package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) RegisterTwitterData(characterId int, twitterId string) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(twitterId, validation.Required)
	if err != nil {
		return errors.ErrValidate
	}

	err = repository.RegisterTwitterData(characterId, twitterId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
