package room

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) PostRoomMessage(characterId int, message *model.RoomPostMessage) error {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	if message == nil {
		return usecaseErrors.ErrValidate
	}

	err := validation.ValidateStruct(message,
		validation.Field(&message.Name, validator.IsNotContainSpecialRune),
		validation.Field(&message.Message, validator.IsNotOnlySpace, validator.IsNotContainSpecialRune),
		validation.Field(&message.ReplyPermission, validation.In("DISALLOW", "FOLLOW", "FOLLOWED", "MUTUAL_FOLLOW", "ALL")),
	)
	if err != nil {
		return usecaseErrors.ErrValidate
	}

	permissions, banned, err := repository.RetrieveRoomOwnPermissions(characterId, message.Room)
	if err != nil {
		return usecaseErrors.ErrValidate
	}
	if banned {
		return errors.New("BANされています")
	}

	if !permissions.Write {
		return errors.New("書き込み権限がありません")
	}

	err = repository.PostRoomMessage(characterId, message, config.GetString("upload-path"))
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
