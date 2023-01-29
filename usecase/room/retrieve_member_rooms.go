package room

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"golang.org/x/sync/errgroup"
)

func (s *RoomUsecase) RetrieveMemberRooms(characterId int) (membereds *[]model.RoomListItem, inviteds *[]model.RoomListItem, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	var eg errgroup.Group

	eg.Go(func() error {
		membereds, err = repository.RetrieveMemberRooms(characterId)
		if err != nil {
			logger.Error(err)
			return err
		}
		return nil
	})

	eg.Go(func() error {
		inviteds, err = repository.RetrieveInvitedRooms(characterId)
		if err != nil {
			logger.Error(err)
			return err
		}
		return nil
	})

	if err = eg.Wait(); err != nil {
		return nil, nil, err
	}

	return membereds, inviteds, nil
}
