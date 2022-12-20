package character

import (
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *CharacterUsecase) SignIn(key string, password string) (id int, notificationToken string, isAdministrator bool, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	if keyInt, err := strconv.Atoi(key); err == nil {
		// ログインキーがID形式だった場合
		id = keyInt
	} else if err := validation.Validate(&key, validation.Required, is.Email); err == nil {
		// ログインキーがメールアドレス形式だった場合
		id, err = repository.ExchangeEmailToId(key)
		if err != nil {
			logger.Error(err)
			return 0, "", false, err
		}
	} else {
		// ログインキーの形式がユーザーIDの可能性があるとき
		id, err = repository.ExchangeUsernameToId(key)
		if err != nil {
			logger.Error(err)
			return 0, "", false, err
		}
	}

	err = validation.Validate(password, validator.IsPassword)
	if err != nil {
		return 0, "", false, usecaseErrors.ErrValidate
	}

	savedPassword, notificationToken, isAdministrator, err := repository.RetrieveCredentials(id)
	if err != nil {
		return 0, "", false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(password))
	if err != nil {
		return 0, "", false, err
	}

	return id, notificationToken, isAdministrator, nil
}
