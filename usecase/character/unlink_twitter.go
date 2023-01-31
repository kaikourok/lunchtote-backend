package character

func (s *CharacterUsecase) UnlinkTwitter(characterId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.UnlinkTwitter(characterId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
