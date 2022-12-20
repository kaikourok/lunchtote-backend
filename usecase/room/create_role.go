package room

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) CreateRole(characterId int, roomId int, roleName string, role *model.RoomRolePermission) (roleId int, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(roomId, validation.Min(1))
	if err != nil {
		return 0, errors.ErrValidate
	}

	err = validation.Validate(roleName, validation.Required, validator.IsNotContainSpecialRune, validator.IsNotOnlySpace)
	if err != nil {
		return 0, errors.ErrValidate
	}

	roleId, err = repository.CreateRole(characterId, roomId, roleName, role)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return roleId, nil
}
