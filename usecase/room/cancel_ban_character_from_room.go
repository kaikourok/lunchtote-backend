package room

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) CancelBanCharacterFromRoom(userId int, targetId int, roomId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(targetId, validation.Min(1))
	if err != nil {
		return errors.ErrValidate
	}

	err = validation.Validate(roomId, validation.Min(1))
	if err != nil {
		return errors.ErrValidate
	}

	err = repository.CancelBanCharacterFromRoom(userId, targetId, roomId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
