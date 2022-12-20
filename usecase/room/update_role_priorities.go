package room

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/library/slice"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) UpdateRolePriorities(characterId int, roomId int, priorities *[]model.RoomRolePriority) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(roomId, validation.Min(1))
	if err != nil {
		return usecaseErrors.ErrValidate
	}

	if priorities == nil {
		return usecaseErrors.ErrValidate
	}

	for _, priority := range *priorities {
		err := validation.ValidateStruct(priority,
			validation.Field(priority.Priority, validation.Min(1)),
			validation.Field(priority.Role, validation.Min(1)),
		)
		if err != nil {
			return usecaseErrors.ErrValidate
		}
	}

	payloadRoles := make([]int, len(*priorities))
	payloadPriorities := make([]int, len(*priorities))

	for i, v := range *priorities {
		payloadRoles[i] = v.Role
		payloadPriorities[i] = v.Priority
	}

	if slice.IsContainsDuplicateValue(&payloadRoles) {
		return usecaseErrors.ErrValidate
	}

	if slice.IsContainsDuplicateValue(&payloadPriorities) {
		return usecaseErrors.ErrValidate
	}

	err = repository.UpdateRolePriorities(characterId, roomId, priorities)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
