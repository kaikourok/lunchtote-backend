package forum

import "github.com/kaikourok/lunchtote-backend/registry"

type ForumUsecase struct {
	registry registry.Registry
}

func NewForumUsecase(registry registry.Registry) *ForumUsecase {
	return &ForumUsecase{registry}
}
