package room

func (s *RoomUsecase) RenameRoomMessageFetchConfig(characterId, configId int, name string) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.RenameRoomMessageFetchConfig(characterId, configId, name)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
