package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) RetrieveFollowList(characterId, targetId int) (list *[]model.CharacterListItem, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(targetId, validation.Min(1))
	if err != nil {
		return nil, errors.ErrValidate
	}

	list, err = repository.RetrieveFollowList(characterId, targetId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return list, nil
}
