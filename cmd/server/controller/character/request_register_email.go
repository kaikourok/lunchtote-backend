package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RequestRegisterEmail(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Email string `json:"email"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.RequestRegisterEmail(session.Get("cid").(int), payload.Email)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
