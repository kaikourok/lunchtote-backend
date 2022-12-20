package control

import (
	"github.com/kaikourok/lunchtote-backend/usecase/control"
)

type ControlController struct {
	usecase *control.ControlUsecase
}

func NewControlController(usecase *control.ControlUsecase) *ControlController {
	return &ControlController{usecase: usecase}
}
