package mail

func (s *MailUsecase) SendAdministratorMail(targetId *int, name string, title string, message string) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.SendAdministratorMail(targetId, name, title, message)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
