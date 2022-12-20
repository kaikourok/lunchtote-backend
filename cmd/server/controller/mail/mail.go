package mail

import (
	"github.com/kaikourok/lunchtote-backend/usecase/mail"
)

type MailController struct {
	usecase *mail.MailUsecase
}

func NewMailController(usecase *mail.MailUsecase) *MailController {
	return &MailController{usecase: usecase}
}
