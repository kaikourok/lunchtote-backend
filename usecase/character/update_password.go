package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	"github.com/kaikourok/lunchtote-backend/library/secure"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *CharacterUsecase) UpdatePassword(characterId int, oldPassword, newPassword string) error {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	err := validation.Validate(oldPassword, validator.IsPassword)
	if err != nil {
		return errors.ErrValidate
	}

	err = validation.Validate(newPassword, validator.IsPassword)
	if err != nil {
		return errors.ErrValidate
	}

	savedPassword, err := repository.RetrievePassword(characterId)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(oldPassword))
	if err != nil {
		logger.Error(err)
		return err
	}

	hashedPassword, err := secure.HashPassword(newPassword, config.GetInt("secure.bcrypt-cost"))
	if err != nil {
		logger.Error(err)
		return err
	}

	err = repository.UpdatePassword(characterId, hashedPassword)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
