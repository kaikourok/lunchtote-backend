package mail

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *MailUsecase) RetrieveSentMails(characterId, start int) (mails *[]model.SentMail, isContinue bool, err error) {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	mailsPerPage := config.GetInt("general.mails-per-page")

	mails, isContinue, err = repository.RetrieveSentMails(characterId, mailsPerPage, start)
	if err != nil {
		logger.Error(err)
		return nil, false, err
	}

	return mails, isContinue, nil
}
