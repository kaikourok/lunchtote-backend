package room

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) RetrieveRoomBanStates(roomId int) (states *[]model.RoomBanState, master int, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(roomId, validation.Min(1))
	if err != nil {
		return nil, 0, errors.ErrValidate
	}

	states, master, err = repository.RetrieveRoomBanStates(roomId)
	if err != nil {
		logger.Error(err)
		return nil, 0, err
	}

	return states, master, nil
}
