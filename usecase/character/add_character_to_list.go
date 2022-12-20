package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) AddCharacterToList(characterId int, targetId int, listId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(targetId, validation.Min(1))
	if err != nil {
		return errors.ErrValidate
	}

	err = validation.Validate(listId, validation.Min(1))
	if err != nil {
		return errors.ErrValidate
	}

	err = repository.AddCharacterToList(characterId, targetId, listId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
