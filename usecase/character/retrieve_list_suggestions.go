package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) RetrieveListSuggestions(characterId int, searchText string, listId int) (suggestions *model.CharacterSuggestions, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(listId, validation.Min(1))
	if err != nil {
		return nil, errors.ErrValidate
	}

	results, err := repository.RetrieveListSuggestions(characterId, searchText, listId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return results.ToDomain(), nil
}
