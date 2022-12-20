package room

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) CreateRoom(characterId int, room *model.Room) (roomId int, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	if room == nil {
		return 0, errors.ErrValidate
	}

	err = validation.ValidateStruct(
		room,
		validation.Field(&room.Title, validation.Required, validator.IsNotOnlySpace, validator.IsNotContainSpecialRune),
		validation.Field(&room.Summary, validator.IsNotContainSpecialRune),
		validation.Field(&room.Description, validator.IsNotContainSpecialRune),
		validation.Field(&room.Tags, validation.Each(validation.Required, validator.IsNotOnlySpace, validator.IsNotContainSpecialRune)),
	)
	if err != nil {
		return 0, errors.ErrValidate
	}

	roomId, err = repository.CreateRoom(characterId, room)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return roomId, nil
}
