package control

import (
	"crypto/sha256"
	"encoding/hex"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	"github.com/kaikourok/lunchtote-backend/library/secure"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *ControlUsecase) CreateDummyCharacters(number int, name string, nickname string, summary string, profile string, password string) error {
	repository := s.registry.GetRepository()
	config := s.registry.GetConfig()
	logger := s.registry.GetLogger()

	err := validation.Validate(name, validation.Required, validator.IsNotContainSpecialRune, validator.IsNotOnlySpace)
	if err != nil {
		return errors.ErrValidate
	}

	err = validation.Validate(name, validation.Required, validator.IsNotContainSpecialRune, validator.IsNotOnlySpace)
	if err != nil {
		return errors.ErrValidate
	}

	err = validation.Validate(summary, validator.IsNotContainSpecialRune)
	if err != nil {
		return errors.ErrValidate
	}

	err = validation.Validate(profile, validator.IsNotContainSpecialRune)
	if err != nil {
		return errors.ErrValidate
	}

	builder := password + config.GetString("secure.frontend-hash-salt")
	for i := 0; i < config.GetInt("secure.frontend-hash-stretch"); i++ {
		b := sha256.Sum256([]byte(builder))
		builder = hex.EncodeToString(b[:])
	}

	hashedPassword, err := secure.HashPassword(builder, config.GetInt("secure.bcrypt-cost"))
	if err != nil {
		logger.Error(err)
		return err
	}

	err = repository.CreateDummyCharacters(number, name, nickname, summary, profile, hashedPassword, func() string {
		return secure.GenerateSecureRandomHex(config.GetInt("secure.notification-token-length"))
	})
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
