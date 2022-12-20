package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) CreateProcessSchema(characterId, groupId int, process *model.CharacterIconProcessSchema) (id int, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	if process == nil {
		return 0, errors.ErrValidate
	}

	err = validation.ValidateStruct(process,
		validation.Field(&process.Name, validation.Required, validator.IsNotOnlySpace, validator.IsNotContainSpecialRune),
	)
	if err != nil {
		return 0, errors.ErrValidate
	}

	id, err = repository.CreateProcessSchema(characterId, groupId, process)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return
}
