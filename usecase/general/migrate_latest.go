package general

func (s *GeneralUsecase) MigrateLatest() error {
	repository := s.registry.GetRepository()
	logger := s.registry.GetLogger()

	err := repository.MigrateLatest()
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
