package room

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *RoomUsecase) UpdateRoomSettings(characterId, roomId int, room *model.Room) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.UpdateRoomSettings(characterId, roomId, room)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
