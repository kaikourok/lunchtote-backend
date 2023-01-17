package room

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *RoomController) UpdateRoomMessageFetchConfigOrders(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Orders *[]model.RoomMessageFetchConfigOrder `json:"orders"`
	}
	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	err = u.usecase.UpdateRoomMessageFetchConfigOrders(session.Get("cid").(int), payload.Orders)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
