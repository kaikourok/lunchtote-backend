package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *CharacterUsecase) UpdateOtherSettings(characterId int, settings *model.CharacterOtherSettings) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.UpdateOtherSettings(characterId, settings)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
