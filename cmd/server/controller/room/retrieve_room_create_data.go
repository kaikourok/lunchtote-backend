package room

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *RoomController) RetrieveRoomCreateData(c *gin.Context) {
	session := sessions.Default(c)

	childrenCreatableRooms, err := u.usecase.RetrieveRoomCreateData(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"childrenCreatableRooms": childrenCreatableRooms,
	})
}
