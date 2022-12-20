package forum

import (
	"github.com/kaikourok/lunchtote-backend/usecase/forum"
)

type ForumController struct {
	usecase *forum.ForumUsecase
}

func NewForumController(usecase *forum.ForumUsecase) *ForumController {
	return &ForumController{usecase: usecase}
}
