package mail

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *MailUsecase) RetrieveMails(characterId, start int, unreadOnly bool) (mails *[]model.ReceivedMail, isContinue bool, err error) {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	mailsPerPage := config.GetInt("general.mails-per-page")

	mails, isContinue, err = repository.RetrieveMails(characterId, unreadOnly, mailsPerPage, start)
	if err != nil {
		logger.Error(err)
		return nil, false, err
	}

	return mails, isContinue, nil
}
