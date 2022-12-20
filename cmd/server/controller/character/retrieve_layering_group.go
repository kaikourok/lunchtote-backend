package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RetrieveLayeringGroup(c *gin.Context) {
	session := sessions.Default(c)

	groupId, err := strconv.Atoi(c.Param("layering_group"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	layeringGroup, err := u.usecase.RetrieveLayeringGroup(session.Get("cid").(int), groupId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, layeringGroup)
}
