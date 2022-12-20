package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *CharacterUsecase) RetrieveLayeringGroupOverviews(characterId int) (overviews *[]model.CharacterIconLayeringGroupOverview, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	overviews, err = repository.RetrieveLayeringGroupOverviews(characterId)
	if err != nil {
		logger.Error(err)
	}

	return
}
