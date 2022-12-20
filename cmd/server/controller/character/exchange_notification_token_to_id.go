package character

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *CharacterController) ExchangeNotificationTokenToId(c *gin.Context) {
	var payload struct {
		Token string `json:"token"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	id, err := u.usecase.ExchangeNotificationTokenToId(payload.Token)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
