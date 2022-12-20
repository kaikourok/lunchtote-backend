package character

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RequestPasswordResetCode(c *gin.Context) {
	var payload struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.RequestPasswordResetCode(payload.Id, payload.Email)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}

	c.Status(http.StatusOK)
}
