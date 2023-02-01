package mail

import (
	"strings"

	"github.com/kaikourok/lunchtote-backend/entity/service"
)

func (s *MailUsecase) SendMail(characterId int, targetId int, title string, message string) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()
	config := s.registry.GetConfig()
	notificator := s.registry.GetNotificator()

	userName, targetWebhook, err := repository.SendMail(characterId, targetId, title, message)
	if err != nil {
		logger.Error(err)
		return err
	}

	if targetWebhook != "" {
		replacer := strings.NewReplacer(
			"{base-path}", config.GetString("general.client-host"),
			"{entry-number-text}", service.ConvertCharacterIdToText(characterId),
			"{name}", userName,
		)

		go notificator.SendWebhook(targetWebhook, replacer.Replace(config.GetString("notification.mail-template")))
	}

	return nil
}
