package character

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *CharacterController) CheckUsernameExists(c *gin.Context) {
	var payload struct {
		Username string `json:"username"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	exists, err := u.usecase.CheckUsernameExists(payload.Username)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exists": exists,
	})
}
