package room

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *RoomController) UpdateRolePriorities(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Priorities *[]model.RoomRolePriority `json:"priorities"`
	}
	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	roomId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.UpdateRolePriorities(session.Get("cid").(int), roomId, payload.Priorities)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
