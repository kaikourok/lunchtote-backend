package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *CharacterUsecase) RetrieveNotifications(characterId int, start int) (notifications []model.Notification, isContinue bool, err error) {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	notifications, isContinue, err = repository.RetrieveNotifications(characterId, start, config.GetInt("general.notifications-per-page"))
	if err != nil {
		logger.Error(err)
		return nil, false, err
	}

	return notifications, isContinue, nil
}
