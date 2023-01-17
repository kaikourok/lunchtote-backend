package room

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *RoomController) DeleteRoomMessageFetchConfig(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Config int `json:"config"`
	}
	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	err = u.usecase.DeleteRoomMessageFetchConfig(session.Get("cid").(int), payload.Config)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
