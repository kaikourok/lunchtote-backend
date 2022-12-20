package room

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *RoomController) CreateRole(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		model.RoomRolePermission
		Name string `json:"name"`
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

	roleId, err := u.usecase.CreateRole(session.Get("cid").(int), roomId, payload.Name, &payload.RoomRolePermission)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": roleId,
	})
}
