package general

func (s *GeneralUsecase) Connections() (connections int, err error) {
	repository := s.registry.GetRepository()
	logger := s.registry.GetLogger()

	connections, err = repository.Connections()
	if err != nil {
		logger.Error(err)
	}

	return
}
