package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) UpdateProfileImages(characterId int, images *[]model.ProfileImage) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(images, validation.Required, validation.Each(validator.IsImagePath(&characterId)))
	if err != nil {
		return errors.ErrValidate
	}

	err = repository.UpdateProfileImages(characterId, images)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
