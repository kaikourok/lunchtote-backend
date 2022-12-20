package mail

import "github.com/kaikourok/lunchtote-backend/registry"

type MailUsecase struct {
	registry registry.Registry
}

func NewMailUsecase(registry registry.Registry) *MailUsecase {
	return &MailUsecase{registry}
}
