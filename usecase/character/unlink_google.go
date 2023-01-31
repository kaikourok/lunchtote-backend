package character

func (s *CharacterUsecase) UnlinkGoogle(characterId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.UnlinkGoogle(characterId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
