package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) ExchangeNotificationTokenToId(token string) (characterId int, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(token, validation.Required, is.Hexadecimal)
	if err != nil {
		return 0, errors.ErrValidate
	}

	characterId, err = repository.ExchangeNotificationTokenToId(token)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return characterId, nil
}
