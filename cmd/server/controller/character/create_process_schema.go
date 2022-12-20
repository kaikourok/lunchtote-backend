package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *CharacterController) CreateProcessSchema(c *gin.Context) {
	session := sessions.Default(c)

	var payload model.CharacterIconProcessSchema
	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	groupId, err := strconv.Atoi(c.Param("layering_group"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	processId, err := u.usecase.CreateProcessSchema(session.Get("cid").(int), groupId, &payload)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": processId,
	})
}
