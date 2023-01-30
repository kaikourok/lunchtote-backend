package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) RetrieveList(characterId, listId int) (listName string, characters []model.CharacterOverview, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(listId, validation.Min(1))
	if err != nil {
		logger.Error(err)
		return "", nil, errors.ErrValidate
	}

	owner, err := repository.RetrieveListOwner(listId)
	if err != nil {
		logger.Error(err)
		return "", nil, err
	}
	if owner != characterId {
		return "", nil, errors.ErrPermission
	}

	listName, characters, err = repository.RetrieveList(listId)
	if err != nil {
		logger.Error(err)
		return "", nil, err
	}

	return
}
