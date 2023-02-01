package mail

import "strings"

func (s *MailUsecase) SendAdministratorMail(targetId *int, name string, title string, message string) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()
	config := s.registry.GetConfig()
	notificator := s.registry.GetNotificator()

	webhooks, err := repository.SendAdministratorMail(targetId, name, title, message)
	if err != nil {
		logger.Error(err)
		return err
	}

	if 0 < len(webhooks) {
		replacer := strings.NewReplacer(
			"{base-path}", config.GetString("general.client-host"),
			"{name}", name,
		)
		notificationMessage := replacer.Replace(config.GetString("notification.administrator-mail-template"))

		go func() {
			for _, webhook := range webhooks {
				notificator.SendWebhook(webhook, notificationMessage)
			}
		}()
	}

	return nil
}
