package room

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) GrantRoomRole(characterId int, targetId int, roleId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(targetId, validation.Min(1))
	if err != nil {
		return errors.ErrValidate
	}

	err = validation.Validate(roleId, validation.Min(1))
	if err != nil {
		return errors.ErrValidate
	}

	err = repository.GrantRoomRole(characterId, targetId, roleId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
