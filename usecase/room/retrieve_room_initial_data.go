package room

import (
	"errors"

	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *RoomUsecase) RetrieveRoomInitialData(characterId, roomId int) (title string, relations *model.RoomRelations, permissions *model.RoomMemberPermission, banned bool, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	type titleResultStruct struct {
		title string
		err   error
	}

	type permissionsResultStruct struct {
		permissions *model.RoomMemberPermission
		banned      bool
		err         error
	}

	type relationsResultStruct struct {
		relations *model.RoomRelations
		err       error
	}

	titleChannel := make(chan titleResultStruct)
	permissionsChannnel := make(chan permissionsResultStruct)
	relationsChannel := make(chan relationsResultStruct)

	go func() {
		title, err := repository.RetrieveRoomTitle(roomId)
		titleChannel <- titleResultStruct{title, err}
	}()

	go func() {
		permissions, _, banned, err := repository.RetrieveRoomOwnPermissions(characterId, roomId)
		permissionsChannnel <- permissionsResultStruct{permissions, banned, err}
	}()

	go func() {
		relations, err := repository.RetrieveRoomRelations(roomId)
		relationsChannel <- relationsResultStruct{relations, err}
	}()

	titleResult := <-titleChannel
	permissionsResult := <-permissionsChannnel
	relationsResult := <-relationsChannel

	if titleResult.err != nil || permissionsResult.err != nil || relationsResult.err != nil {
		if titleResult.err != nil {
			logger.Error(titleResult.err)
		}
		if permissionsResult.err != nil {
			logger.Error(permissionsResult.err)
		}
		if relationsResult.err != nil {
			logger.Error(relationsResult.err)
		}
		return "", nil, nil, false, errors.New("データの取得中にエラーが発生しました")
	}

	return titleResult.title, relationsResult.relations, permissionsResult.permissions, permissionsResult.banned, nil
}
