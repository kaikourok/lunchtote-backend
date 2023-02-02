package room

import (
	validation "github.com/go-ozzo/ozzo-validation"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) UnsubscribeRoomMessage(characterId int, roomId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(roomId, validation.Min(1))
	if err != nil {
		return usecaseErrors.ErrValidate
	}

	err = repository.UnsubscribeRoomMessage(characterId, roomId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
