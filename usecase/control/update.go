package control

func (s *ControlUsecase) Update() error {
	repository := s.registry.GetRepository()
	logger := s.registry.GetLogger()

	err := repository.Update()
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
