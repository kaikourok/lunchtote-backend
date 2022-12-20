package mail

func (s *MailUsecase) SendMail(characterId int, targetId int, title string, message string) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.SendMail(characterId, targetId, title, message)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
