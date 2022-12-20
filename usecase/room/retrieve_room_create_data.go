package room

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *RoomUsecase) RetrieveRoomCreateData(characterId int) (childrenCreatableRooms *[]model.RoomOverview, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	rooms, err := repository.RetrieveChildrenCreatableRooms(characterId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return rooms, nil
}
