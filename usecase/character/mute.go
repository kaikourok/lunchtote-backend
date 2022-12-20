package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) Mute(characterId int, targetId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(targetId, validation.Min(1))
	if err != nil {
		return errors.ErrValidate
	}

	err = repository.Mute(characterId, targetId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
