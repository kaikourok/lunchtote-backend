package room

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *RoomUsecase) RetrieveOwnedRooms(characterId int) (rooms *[]model.RoomListItem, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	rooms, err = repository.RetrieveOwnedRooms(characterId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return rooms, nil
}
