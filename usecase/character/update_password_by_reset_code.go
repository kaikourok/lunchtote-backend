package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	"github.com/kaikourok/lunchtote-backend/library/secure"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) UpdatePasswordByResetCode(characterId int, code, newPassword string) error {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	err := validation.Validate(code, validation.Required, is.Hexadecimal)
	if err != nil {
		return errors.ErrValidate
	}

	err = validation.Validate(newPassword, validator.IsPassword)
	if err != nil {
		return errors.ErrValidate
	}

	hashedPassword, err := secure.HashPassword(newPassword, config.GetInt("secure.bcrypt-cost"))
	if err != nil {
		logger.Error(err)
		return err
	}

	return repository.UpdatePasswordByResetCode(characterId, code, hashedPassword)
}
