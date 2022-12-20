package room

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *RoomController) DeleteRole(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Role int `json:"role"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	err = u.usecase.DeleteRole(session.Get("cid").(int), payload.Role)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
