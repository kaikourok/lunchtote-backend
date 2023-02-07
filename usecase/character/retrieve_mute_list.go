package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *CharacterUsecase) RetrieveMuteList(characterId int) (list *[]model.GeneralCharacterListItem, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	list, err = repository.RetrieveMuteList(characterId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return list, nil
}
