package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) RenameList(characterId int, listId int, newName string) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(listId, validation.Min(1))
	if err != nil {
		return errors.ErrValidate
	}

	err = validation.Validate(newName, validation.Required, validator.IsNotContainSpecialRune, validator.IsNotOnlySpace)
	if err != nil {
		return errors.ErrValidate
	}

	owner, err := repository.RetrieveListOwner(listId)
	if err != nil {
		logger.Error(err)
		return errors.ErrValidate
	}
	if owner != characterId {
		return errors.ErrPermission
	}

	err = repository.RenameList(listId, newName)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
