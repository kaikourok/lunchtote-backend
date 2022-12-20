package character

import "github.com/kaikourok/lunchtote-backend/registry"

type CharacterUsecase struct {
	registry registry.Registry
}

func NewCharacterUsecase(registry registry.Registry) *CharacterUsecase {
	return &CharacterUsecase{registry}
}
