package control

import "github.com/kaikourok/lunchtote-backend/registry"

type ControlUsecase struct {
	registry registry.Registry
}

func NewControlUsecase(registry registry.Registry) *ControlUsecase {
	return &ControlUsecase{registry}
}
