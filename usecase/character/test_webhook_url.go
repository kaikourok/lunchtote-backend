package character

func (s *CharacterUsecase) TestWebhookUrl(url string) error {
	notificator := s.registry.GetNotificator()

	go notificator.SendWebhook(url, "Webhookのテスト送信です")

	return nil
}
