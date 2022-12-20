package character

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/service"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) Follow(characterId int, targetId int) error {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()
	notificator := s.registry.GetNotificator()

	err := validation.Validate(targetId, validation.Min(1))
	if err != nil {
		return errors.ErrValidate
	}

	userName, webhook, err := repository.Follow(characterId, targetId)
	if err != nil {
		logger.Error(err)
		return err
	}

	defer func() {
		if webhook != "" {
			s := strings.ReplaceAll(
				strings.ReplaceAll(
					config.GetString("notification.followed-template"),
					"{entry-number}", service.ConvertCharacterIdToText(characterId),
				),
				"{name}", userName,
			)
			go notificator.SendWebhook(webhook, s)
		}
	}()

	return nil
}
