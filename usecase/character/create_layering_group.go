package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) CreateLayeringGroup(characterId int, name string) (id int, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(name, validation.Required, validator.IsNotOnlySpace, validator.IsNotContainSpecialRune)
	if err != nil {
		return 0, errors.ErrValidate
	}

	id, err = repository.CreateLayeringGroup(characterId, name)
	if err != nil {
		logger.Error(err)
	}

	return
}
