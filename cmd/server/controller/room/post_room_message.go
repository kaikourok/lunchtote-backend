package room

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *RoomController) PostRoomMessage(c *gin.Context) {
	session := sessions.Default(c)

	var payload model.RoomPostMessage
	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	err = u.usecase.PostRoomMessage(session.Get("cid").(int), &payload)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
