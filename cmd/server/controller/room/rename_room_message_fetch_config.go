package room

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *RoomController) RenameRoomMessageFetchConfig(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Config int    `json:"config"`
		Name   string `json:"name"`
	}
	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	err = u.usecase.RenameRoomMessageFetchConfig(session.Get("cid").(int), payload.Config, payload.Name)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
