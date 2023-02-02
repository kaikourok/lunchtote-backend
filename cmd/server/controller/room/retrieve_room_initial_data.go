package room

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *RoomController) RetrieveRoomInitialData(c *gin.Context) {
	session := sessions.Default(c)

	roomId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	initialData, err := u.usecase.RetrieveRoomInitialData(session.Get("cid").(int), roomId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, initialData)
}
