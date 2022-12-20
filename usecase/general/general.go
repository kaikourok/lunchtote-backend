package general

import "github.com/kaikourok/lunchtote-backend/registry"

type GeneralUsecase struct {
	registry registry.Registry
}

func NewGeneralUsecase(registry registry.Registry) *GeneralUsecase {
	return &GeneralUsecase{registry}
}
