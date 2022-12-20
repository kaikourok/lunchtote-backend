package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) DeleteProcessSchema(c *gin.Context) {
	session := sessions.Default(c)

	processId, err := strconv.Atoi(c.Param("process"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.DeleteProcessSchema(session.Get("cid").(int), processId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
