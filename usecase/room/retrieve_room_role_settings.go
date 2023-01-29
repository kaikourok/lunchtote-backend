package room

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) RetrieveRoomRoleSettings(characterId, roomId int) (roles []model.RoomRole, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(roomId, validation.Min(1))
	if err != nil {
		return nil, usecaseErrors.ErrValidate
	}

	roles, master, err := repository.RetrieveRoomRoleSettings(roomId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if characterId != master {
		err := errors.New("ルームの管理権限を持っていません")
		return nil, err
	}

	return roles, nil
}
