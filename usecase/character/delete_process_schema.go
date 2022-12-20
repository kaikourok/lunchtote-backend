package character

func (s *CharacterUsecase) DeleteProcessSchema(characterId, processId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.DeleteProcessSchema(characterId, processId)
	if err != nil {
		logger.Error(err)
	}

	return err
}
