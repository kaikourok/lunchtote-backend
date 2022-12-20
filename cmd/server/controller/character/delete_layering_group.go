package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) DeleteLayeringGroup(c *gin.Context) {
	session := sessions.Default(c)

	layeringGroup, err := strconv.Atoi(c.Param("layering_group"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.DeleteLayeringGroup(session.Get("cid").(int), layeringGroup)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
