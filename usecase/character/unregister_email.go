package character

func (s *CharacterUsecase) UnregisterEmail(characterId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.UnregisterEmail(characterId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
