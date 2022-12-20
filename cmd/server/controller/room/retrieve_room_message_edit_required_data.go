package room

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *RoomController) RetrieveRoomMessageEditRequiredData(c *gin.Context) {
	session := sessions.Default(c)

	data, err := u.usecase.RetrieveRoomMessageEditRequiredData(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}
