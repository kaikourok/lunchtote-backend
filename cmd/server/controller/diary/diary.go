package diary

import (
	"github.com/kaikourok/lunchtote-backend/usecase/diary"
)

type DiaryController struct {
	usecase *diary.DiaryUsecase
}

func NewDiaryController(usecase *diary.DiaryUsecase) *DiaryController {
	return &DiaryController{usecase: usecase}
}
