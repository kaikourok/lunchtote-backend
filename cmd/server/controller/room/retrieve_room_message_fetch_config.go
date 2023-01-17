package room

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *RoomController) RetrieveRoomMessageFetchConfig(c *gin.Context) {
	session := sessions.Default(c)

	configs, err := u.usecase.RetrieveRoomMessageFetchConfig(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, configs)
}
