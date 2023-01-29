package room

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) RetrieveRoomBanStates(characterId, roomId int) (states *[]model.RoomBanState, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(roomId, validation.Min(1))
	if err != nil {
		return nil, errors.ErrValidate
	}

	permissions, _, banned, err := repository.RetrieveRoomOwnPermissions(characterId, roomId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if banned || !permissions.Ban {
		return nil, errors.ErrPermission
	}

	states, err = repository.RetrieveRoomBanStates(roomId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return states, nil
}
