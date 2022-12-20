package character

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *CharacterController) SignUp(c *gin.Context) {
	var payload struct {
		Name     string  `json:"name"`
		Nickname string  `json:"nickname"`
		Username string  `json:"username"`
		Password string  `json:"password"`
		Email    *string `json:"email"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	id, err := u.usecase.SignUp(payload.Name, payload.Nickname, payload.Username, payload.Password, payload.Email)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
