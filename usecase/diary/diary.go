package diary

import "github.com/kaikourok/lunchtote-backend/registry"

type DiaryUsecase struct {
	registry registry.Registry
}

func NewDiaryUsecase(registry registry.Registry) *DiaryUsecase {
	return &DiaryUsecase{registry}
}
