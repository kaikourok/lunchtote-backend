package room

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *RoomController) RetrieveRoomInviteSuggestions(c *gin.Context) {
	session := sessions.Default(c)

	roomId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var payload struct {
		Text string `json:"text"`
	}

	err = c.BindJSON(&payload)
	if err != nil {
		return
	}

	results, err := u.usecase.RetrieveRoomInviteSuggestions(session.Get("cid").(int), payload.Text, roomId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, results)
}
