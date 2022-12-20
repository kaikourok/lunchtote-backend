package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *CharacterUsecase) RetrieveBlockList(characterId int) (list *[]model.CharacterListItem, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	list, err = repository.RetrieveBlockList(characterId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return list, nil
}
