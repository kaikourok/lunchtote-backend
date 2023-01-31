package character

func (s *CharacterUsecase) DeleteCharacter(characterId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.DeleteCharacter(characterId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
