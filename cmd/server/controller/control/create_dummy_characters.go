package control

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *ControlController) CreateDummyCharacters(c *gin.Context) {
	var payload struct {
		Number   int    `json:"number"`
		Name     string `json:"name"`
		Nickname string `json:"nickname"`
		Summary  string `json:"summary"`
		Profile  string `json:"profile"`
		Password string `json:"password"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.CreateDummyCharacters(payload.Number, payload.Name, payload.Nickname, payload.Nickname, payload.Profile, payload.Password)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
