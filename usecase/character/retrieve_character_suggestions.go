package character

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *CharacterUsecase) RetrieveCharacterSuggestions(characterId int, searchText string, excludeOwn bool) (suggestions *model.CharacterSuggestions, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	results, err := repository.RetrieveCharacterSuggestions(characterId, searchText, excludeOwn)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return results.ToDomain(), nil
}
