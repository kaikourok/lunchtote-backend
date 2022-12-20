package room

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *RoomController) RetrieveRoomBanStates(c *gin.Context) {
	session := sessions.Default(c)

	roomId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	states, master, err := u.usecase.RetrieveRoomBanStates(roomId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if master != session.Get("cid").(int) {
		c.Status(http.StatusForbidden)
		return
	}

	c.JSON(http.StatusOK, states)
}
