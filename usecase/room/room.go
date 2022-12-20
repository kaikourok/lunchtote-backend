package room

import "github.com/kaikourok/lunchtote-backend/registry"

type RoomUsecase struct {
	registry registry.Registry
}

func NewRoomUsecase(registry registry.Registry) *RoomUsecase {
	return &RoomUsecase{registry}
}
