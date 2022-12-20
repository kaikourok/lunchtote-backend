package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *CharacterUsecase) RetrieveCharacterList(characterId *int, page int) (list *[]model.CharacterListItem, maxId int, err error) {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	characterListItemPerPage := config.GetInt("general.character-list-items-per-page")

	characterListData, maxId, err := repository.RetrieveCharacterList(characterId, page*characterListItemPerPage+1, (page+1)*characterListItemPerPage)
	if err != nil {
		logger.Error(err)
		return nil, 0, err
	}

	// TODO: MAXIDではなくLAST PAGEで返すべきでは
	return characterListData, maxId, nil
}
