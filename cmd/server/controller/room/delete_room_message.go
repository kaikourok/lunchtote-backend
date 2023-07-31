package room

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *RoomController) DeleteRoomMessage(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Target int `json:"target"`
	}
	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	err = u.usecase.DeleteRoomMessage(session.Get("cid").(int), payload.Target)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
