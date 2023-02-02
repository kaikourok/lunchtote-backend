package general

func (s *GeneralUsecase) MigrateDrop() error {
	repository := s.registry.GetRepository()
	logger := s.registry.GetLogger()

	err := repository.MigrateDrop()
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
