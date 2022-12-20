package general

import (
	"github.com/kaikourok/lunchtote-backend/usecase/general"
)

type GeneralController struct {
	usecase *general.GeneralUsecase
}

func NewGeneralController(usecase *general.GeneralUsecase) *GeneralController {
	return &GeneralController{usecase: usecase}
}
