package character

import (
	"strconv"

	"github.com/kaikourok/lunchtote-backend/library/secure"
	templateLib "github.com/kaikourok/lunchtote-backend/library/template"
)

func (s *CharacterUsecase) RequestRegisterEmail(characterId int, email string) error {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()
	emailSender := s.registry.GetEmail()

	confirmCode := secure.GenerateSecureRandomHex(config.GetInt("email.confirm-code-length"))

	err := repository.SetConfirmCode(characterId, email, confirmCode, config.GetInt("email.confirm-expire"))
	if err != nil {
		logger.Error(err)
		return err
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
		return err
	}

	title, err := templateLib.TemplateExecute(confirmMailTitleTemplate, map[string]any{})
	if err != nil {
		logger.Error(err)
		return err
	}

	err = emailSender.SendEmail(email, title, body)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
