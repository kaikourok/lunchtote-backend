package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) UpdateIcons(characterId int, icons *[]model.Icon, insertOnly bool) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	if icons == nil {
		return errors.ErrValidate
	}

	for _, icon := range *icons {
		err := validation.ValidateStruct(&icon,
			validation.Field(&icon.Path, validator.IsImagePath(&characterId)),
		)
		if err != nil {
			return errors.ErrValidate
		}
	}

	err := repository.UpdateIcons(characterId, icons, insertOnly)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
