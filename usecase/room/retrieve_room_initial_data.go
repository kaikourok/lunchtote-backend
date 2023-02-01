package room

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"golang.org/x/sync/errgroup"
)

func (s *RoomUsecase) RetrieveRoomInitialData(characterId, roomId int) (initialData *model.RoomInitialData, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	initialData = &model.RoomInitialData{}

	var eg errgroup.Group

	eg.Go(func() error {
		title, err := repository.RetrieveRoomTitle(roomId)
		if err != nil {
			logger.Error(err)
			return err
		}

		initialData.Title = title
		return nil
	})

	eg.Go(func() error {
		permissions, _, banned, err := repository.RetrieveRoomOwnPermissions(characterId, roomId)
		if err != nil {
			logger.Error(err)
			return err
		}

		initialData.Permissions = *permissions
		initialData.Banned = banned
		return nil
	})

	eg.Go(func() error {
		relations, err := repository.RetrieveRoomRelations(roomId)
		if err != nil {
			logger.Error(err)
			return err
		}

		initialData.Relations = *relations
		return nil
	})

	eg.Go(func() error {
		states, err := repository.RetrieveRoomSubscribeStates(characterId, roomId)
		if err != nil {
			logger.Error(err)
			return err
		}

		initialData.SubscribeStates = *states
		return nil
	})

	err = eg.Wait()
	return
}
