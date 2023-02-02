package character

import (
	"strconv"
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

	userName, targetWebhook, err := repository.Follow(characterId, targetId)
	if err != nil {
		logger.Error(err)
		return err
	}

	if targetWebhook != "" {
		replacer := strings.NewReplacer(
			"{base-path}", config.GetString("general.client-host"),
			"{entry-number-text}", service.ConvertCharacterIdToText(characterId),
			"{entry-number}", strconv.Itoa(characterId),
			"{name}", userName,
		)

		go notificator.SendWebhook(targetWebhook, replacer.Replace(config.GetString("notification.followed-template")))
	}

	return nil
}
