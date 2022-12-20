package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) RetrieveRelatedFollowerList(characterId, targetId int) (list *[]model.CharacterListItem, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(targetId, validation.Min(1))
	if err != nil {
		return nil, usecaseErrors.ErrValidate
	}

	list, err = repository.RetrieveRelatedFollowerList(characterId, targetId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return list, nil
}
