package character

import (
	"errors"
	"text/template"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	templateLib "github.com/kaikourok/lunchtote-backend/library/template"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

var confirmedMailTitleTemplate, confirmedMailBodyTemplate *template.Template

func init() {
	confirmedMailTitleTemplate = template.Must(template.ParseFiles("template/mail/confirmed/title.gotmpl"))
	confirmedMailBodyTemplate = template.Must(template.ParseFiles("template/mail/confirmed/body.gotmpl"))
}

func (s *CharacterUsecase) ConfirmEmail(characterId int, code string) error {
	repository := s.registry.GetRepository()
	logger := s.registry.GetLogger()
	mailSender := s.registry.GetEmail()

	err := validation.Validate(code, is.Hexadecimal)
	if err != nil {
		return usecaseErrors.ErrValidate
	}

	savedCode, email, err := repository.RetrieveConfirmCode(characterId)
	if err != nil {
		logger.Error(err)
		return err
	}

	if code != savedCode {
		err := errors.New("コードが一致しません")
		logger.Error(err)
		return err
	}

	err = repository.UpdateEmail(characterId, email)
	if err != nil {
		logger.Error(err)
		return err
	}

	defer func() {
		body, err := templateLib.TemplateExecute(confirmedMailBodyTemplate, map[string]any{})
		if err != nil {
			logger.Error(err)
			return
		}

		title, err := templateLib.TemplateExecute(confirmedMailTitleTemplate, map[string]any{})
		if err != nil {
			logger.Error(err)
			return
		}

		err = mailSender.SendEmail(email, title, body)
		if err != nil {
			logger.Error(err)
			return
		}
	}()

	return nil
}
