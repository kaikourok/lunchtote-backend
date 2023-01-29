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

	states, err := u.usecase.RetrieveRoomBanStates(session.Get("cid").(int), roomId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, states)
}
