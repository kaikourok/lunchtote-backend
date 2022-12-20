package mail

func (s *MailUsecase) SetMailRead(characterId int, mailId int) (existsUnreadMail bool, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	existsUnreadMail, err = repository.SetMailRead(characterId, mailId)
	if err != nil {
		logger.Error(err)
		return false, err
	}

	return existsUnreadMail, nil
}
