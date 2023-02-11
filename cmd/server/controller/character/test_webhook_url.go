package character

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *CharacterController) TestWebhookUrl(c *gin.Context) {
	var payload struct {
		Url string `json:"url"`
	}
	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.TestWebhookUrl(payload.Url)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
