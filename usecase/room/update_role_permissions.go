package room

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) UpdateRolePermissions(characterId, roleId int, roleName string, role *model.RoomRolePermission) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(roleId, validation.Min(1))
	if err != nil {
		return usecaseErrors.ErrValidate
	}

	err = validation.Validate(roleName, validation.Required, validator.IsNotContainSpecialRune, validator.IsNotOnlySpace)
	if err != nil {
		return usecaseErrors.ErrValidate
	}

	if role == nil {
		return usecaseErrors.ErrValidate
	}

	err = repository.UpdateRolePermissions(characterId, roleId, roleName, role)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
