package diary

func (s *DiaryUsecase) ClearReservedDiary(characterId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.ClearReservedDiary(characterId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
