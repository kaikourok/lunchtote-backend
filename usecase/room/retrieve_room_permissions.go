package room

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *RoomUsecase) RetrieveRoomOwnPermissions(characterId, roomId int) (permissions *model.RoomMemberPermission, banned bool, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	permissions, _, banned, err = repository.RetrieveRoomOwnPermissions(characterId, roomId)
	if err != nil {
		logger.Error(err)
	}

	return
}
