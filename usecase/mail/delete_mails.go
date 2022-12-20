package mail

func (s *MailUsecase) DeleteMails(characterId int, mailIds *[]int) (deletedIds *[]int, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	deletedIds, err = repository.DeleteMails(characterId, mailIds)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return deletedIds, nil
}
