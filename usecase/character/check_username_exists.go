package character

func (s *CharacterUsecase) CheckUsernameExists(username string) (exists bool, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	exists, err = repository.CheckUsernameExists(username)
	if err != nil {
		logger.Error(err)
	}

	return
}
