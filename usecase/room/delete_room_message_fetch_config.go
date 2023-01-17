package room

func (s *RoomUsecase) DeleteRoomMessageFetchConfig(characterId, configId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.DeleteRoomMessageFetchConfig(characterId, configId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
