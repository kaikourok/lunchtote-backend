package character

import (
	"strconv"
	"text/template"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	"github.com/kaikourok/lunchtote-backend/library/secure"
	templateLib "github.com/kaikourok/lunchtote-backend/library/template"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

var confirmMailTitleTemplate, confirmMailBodyTemplate *template.Template

func init() {
	confirmMailTitleTemplate = template.Must(template.ParseFiles("template/mail/confirm/title.gotmpl"))
	confirmMailBodyTemplate = template.Must(template.ParseFiles("template/mail/confirm/body.gotmpl"))
}

func (s *CharacterUsecase) SignUp(name, nickname, username, password string, email *string) (id int, err error) {
	repository := s.registry.GetRepository()
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	emailSender := s.registry.GetEmail()

	err = validation.Validate(name, validation.Required, validator.IsNotContainSpecialRune, validator.IsNotOnlySpace)
	if err != nil {
		return 0, usecaseErrors.ErrValidate
	}

	err = validation.Validate(nickname, validation.Required, validator.IsNotContainSpecialRune, validator.IsNotOnlySpace)
	if err != nil {
		return 0, usecaseErrors.ErrValidate
	}

	err = validation.Validate(username, validator.IsUsername(config.GetInt("validation.username-min-length"), config.GetInt("validation.username-max-length")))
	if err != nil {
		return 0, usecaseErrors.ErrValidate
	}

	err = validation.Validate(password, validator.IsPassword)
	if err != nil {
		return 0, usecaseErrors.ErrValidate
	}

	if email != nil {
		err = validation.Validate(*email, is.Email)
		if err != nil {
			return 0, usecaseErrors.ErrValidate
		}
	}

	cryptedPassword, err := secure.HashPassword(password, config.GetInt("secure.bcrypt-cost"))
	if err != nil {
		return 0, err
	}

	notificationToken := secure.GenerateSecureRandomHex(config.GetInt("secure.notification-token-length"))
	id, err = repository.CreateCharacter(name, nickname, username, cryptedPassword, notificationToken)
	if err != nil {
		return 0, err
	}

	if email != nil {
		confirmCode := secure.GenerateSecureRandomHex(config.GetInt("email.confirm-code-length"))

		err = repository.SetConfirmCode(id, *email, confirmCode, config.GetInt("email.confirm-expire"))
		if err != nil {
			logger.Error(err)
			return
		}

		confirmExpire := config.GetInt("email.confirm-expire")

		var confirmExpireString string
		if confirmExpire%60 == 0 {
			confirmExpireString = strconv.Itoa(confirmExpire/60) + "時間"
		} else {
			confirmExpireString = strconv.Itoa(confirmExpire) + "分"
		}

		body, err := templateLib.TemplateExecute(confirmMailBodyTemplate, map[string]any{
			"url":    config.GetString("general.host") + config.GetString("email.confirm-path") + confirmCode,
			"expire": confirmExpireString,
		})
		if err != nil {
			logger.Error(err)
			return 0, err
		}

		title, err := templateLib.TemplateExecute(confirmMailTitleTemplate, map[string]any{})
		if err != nil {
			logger.Error(err)
			return 0, err
		}

		err = emailSender.SendEmail(*email, title, body)
		if err != nil {
			logger.Error(err)
			return 0, err
		}
	}

	return id, nil
}
