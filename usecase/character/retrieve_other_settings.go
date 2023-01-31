package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *CharacterUsecase) RetrieveOtherSettings(characterId int) (settings *model.CharacterOtherSettingsState, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	settings, err = repository.RetrieveOtherSettings(characterId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return
}
