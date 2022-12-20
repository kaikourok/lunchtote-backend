package character

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *CharacterUsecase) RetrieveCharacterIconsEditData(characterId int) (icons *[]model.Icon, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	icons, err = repository.RetrieveCharacterIcons(characterId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return icons, nil
}
