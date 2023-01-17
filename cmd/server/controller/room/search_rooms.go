package room

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *RoomController) SearchRooms(c *gin.Context) {
	session := sessions.Default(c)

	var payload model.RoomSearchOptions
	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	rooms, isContinue, err := u.usecase.SearchRooms(session.Get("cid").(int), &payload)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rooms":      rooms,
		"isContinue": isContinue,
	})
}
