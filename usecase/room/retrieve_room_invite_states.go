package room

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) RetrieveRoomInviteStates(roomId int) (states *[]model.RoomInviteState, master int, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(roomId, validation.Min(1))
	if err != nil {
		return nil, 0, usecaseErrors.ErrValidate
	}

	states, master, err = repository.RetrieveRoomInviteStates(roomId)
	if err != nil {
		logger.Error(err)
		return nil, 0, err
	}

	return states, master, nil
}
