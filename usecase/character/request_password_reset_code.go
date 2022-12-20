package character

import (
	"strconv"
	"text/template"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/kaikourok/lunchtote-backend/library/secure"
	templateLib "github.com/kaikourok/lunchtote-backend/library/template"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

var passwordResetMailTitleTemplate, passwordResetMailBodyTemplate *template.Template

func init() {
	passwordResetMailTitleTemplate = template.Must(template.ParseFiles("template/mail/password-reset/title.gotmpl"))
	passwordResetMailBodyTemplate = template.Must(template.ParseFiles("template/mail/password-reset/body.gotmpl"))
}

func (s *CharacterUsecase) RequestPasswordResetCode(characterId int, email string) error {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()
	mailSender := s.registry.GetEmail()

	err := validation.Validate(email, validation.Required, is.Email)
	if err != nil {
		return errors.ErrValidate
	}

	resetExpire := config.GetInt("secure.password-reset-expire")
	resetCode := secure.GenerateSecureRandomHex(config.GetInt("secure.password-reset-code-length"))

	err = repository.SetPasswordResetCode(characterId, email, resetCode, resetExpire)
	if err != nil {
		logger.Error(err)
		return err
	}

	var resetExpireString string
	if resetExpire%60 == 0 {
		resetExpireString = strconv.Itoa(resetExpire/60) + "時間"
	} else {
		resetExpireString = strconv.Itoa(resetExpire) + "分"
	}

	body, err := templateLib.TemplateExecute(passwordResetMailBodyTemplate, map[string]any{
		"url":    config.GetString("general.host") + config.GetString("secure.password-reset-path") + resetCode,
		"expire": resetExpireString,
	})
	if err != nil {
		logger.Error(err)
		return err
	}

	title, err := templateLib.TemplateExecute(passwordResetMailTitleTemplate, map[string]any{})
	if err != nil {
		logger.Error(err)
		return err
	}

	err = mailSender.SendEmail(email, title, body)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
