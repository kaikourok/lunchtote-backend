package room

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) RetrieveRoomMembers(characterId, roomId int) (members *[]model.RomeMemberWithRoles, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(roomId, validation.Min(1))
	if err != nil {
		return nil, usecaseErrors.ErrValidate
	}

	members, err = repository.RetrieveRoomMembers(characterId, roomId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return members, nil
}
