package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) CreateList(characterId int, name string) (listId int, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(name, validation.Required, validator.IsNotOnlySpace, validator.IsNotContainSpecialRune)
	if err != nil {
		return 0, usecaseErrors.ErrValidate
	}

	listId, err = repository.CreateList(characterId, name)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return listId, nil
}
