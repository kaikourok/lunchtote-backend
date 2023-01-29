package room

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *RoomController) UpdateRolePermissions(c *gin.Context) {
	session := sessions.Default(c)

	roleId, err := strconv.Atoi(c.Param("role"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var payload struct {
		Name string `json:"name"`
		model.RoomRolePermission
	}
	err = c.BindJSON(&payload)
	if err != nil {
		return
	}

	err = u.usecase.UpdateRolePermissions(session.Get("cid").(int), roleId, payload.Name, &payload.RoomRolePermission)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
