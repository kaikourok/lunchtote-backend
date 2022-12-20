package room

import (
	"github.com/kaikourok/lunchtote-backend/usecase/room"
)

type RoomController struct {
	usecase *room.RoomUsecase
}

func NewRoomController(usecase *room.RoomUsecase) *RoomController {
	return &RoomController{usecase: usecase}
}
