package room

import (
	validation "github.com/go-ozzo/ozzo-validation"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) RevokeRoomRole(characterId int, targetId int, roleId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(targetId, validation.Min(1))
	if err != nil {
		return usecaseErrors.ErrValidate
	}

	err = validation.Validate(roleId, validation.Min(1))
	if err != nil {
		return usecaseErrors.ErrValidate
	}

	err = repository.RevokeRoomRole(characterId, targetId, roleId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
