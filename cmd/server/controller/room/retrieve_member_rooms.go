package room

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *RoomController) RetrieveMemberRooms(c *gin.Context) {
	session := sessions.Default(c)

	membereds, inviteds, err := u.usecase.RetrieveMemberRooms(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"membereds": membereds,
		"inviteds":  inviteds,
	})
}
