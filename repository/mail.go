package repository

import "github.com/kaikourok/lunchtote-backend/entity/model"

type mailRepository interface {
	RetrieveMails(characterId int, unreadedOnly bool, limit int, start int) (mails *[]model.ReceivedMail, isContinue bool, err error)
	RetrieveSentMails(characterId int, limit int, start int) (mails *[]model.SentMail, isContinue bool, err error)
	SendAdministratorMail(targetId *int, name, title, message string) error
	SendMail(userId, targetId int, title, message string) error
	DeleteMails(characterId int, mailIds *[]int) (deletedIds *[]int, err error)
	SetMailRead(characterId int, mailId int) (existsUnreadMail bool, err error)
}
