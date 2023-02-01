package character

func (s *CharacterUsecase) UpdateNotificationChecked(characterId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.UpdateNotificationChecked(characterId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
