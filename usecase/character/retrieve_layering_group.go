package character

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *CharacterUsecase) RetrieveLayeringGroup(characterId, groupId int) (layeringGroup *model.CharacterIconLayeringGroup, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	layeringGroup, err = repository.RetrieveLayeringGroup(characterId, groupId)
	if err != nil {
		logger.Error(err)
	}

	return
}
