package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) ConfirmEmail(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Code string `json:"code"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.ConfirmEmail(session.Get("cid").(int), payload.Code)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}

	c.Status(http.StatusOK)
}
