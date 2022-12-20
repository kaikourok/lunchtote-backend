package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) UpdateProfile(characterId int, profile *model.ProfileEditData) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(profile, validation.Required)
	if err != nil {
		return errors.ErrValidate
	}

	err = validation.ValidateStruct(profile,
		validation.Field(profile.Name, validation.Required, validator.IsNotContainSpecialRune, validator.IsNotOnlySpace),
		validation.Field(profile.Nickname, validation.Required, validator.IsNotContainSpecialRune, validator.IsNotOnlySpace),
		validation.Field(profile.Summary, validator.IsNotContainSpecialRune),
		validation.Field(profile.Profile, validator.IsNotContainSpecialRune),
		validation.Field(profile.Mainicon, validator.IsImagePath(&characterId)),
		validation.Field(profile.Tags, validation.Each(validator.IsNotOnlySpace, validator.IsNotContainSpecialRune)),
	)
	if err != nil {
		return errors.ErrValidate
	}

	err = repository.UpdateProfile(characterId, profile)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
