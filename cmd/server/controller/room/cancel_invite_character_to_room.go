package room

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *RoomController) CancelInviteCharacterToRoom(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Target int `json:"target"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	roomId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.CancelInviteCharacterToRoom(session.Get("cid").(int), payload.Target, roomId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
