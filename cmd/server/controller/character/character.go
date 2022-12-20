package character

import (
	"github.com/kaikourok/lunchtote-backend/registry"
	"github.com/kaikourok/lunchtote-backend/usecase/character"
)

type CharacterController struct {
	usecase  *character.CharacterUsecase
	registry registry.Registry
}

func NewCharacterController(usecase *character.CharacterUsecase, registry registry.Registry) *CharacterController {
	return &CharacterController{usecase: usecase, registry: registry}
}
