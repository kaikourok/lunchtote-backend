package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *CharacterUsecase) RetrieveLists(characterId int) (lists *[]model.ListOverview, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	lists, err = repository.RetrieveLists(characterId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return lists, nil
}
