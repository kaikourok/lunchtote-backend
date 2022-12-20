package room

func (s *RoomUsecase) DeleteRole(characterId int, roleId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.DeleteRole(characterId, roleId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
