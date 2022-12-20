package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) UpdateLayeringGroupName(characterId, groupId int, name string) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(groupId, validation.Min(1))
	if err != nil {
		return errors.ErrValidate
	}

	err = validation.Validate(name, validation.Required, validator.IsNotOnlySpace, validator.IsNotContainSpecialRune)
	if err != nil {
		return errors.ErrValidate
	}

	err = repository.UpdateLayeringGroupName(characterId, groupId, name)
	if err != nil {
		logger.Error(err)
	}

	return err
}
