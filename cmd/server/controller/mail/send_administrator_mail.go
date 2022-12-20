package mail

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *MailController) SendAdministratorMail(c *gin.Context) {
	var payload struct {
		Target  *int   `json:"target"`
		Name    string `json:"name"`
		Title   string `json:"title"`
		Message string `json:"message"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.SendAdministratorMail(payload.Target, payload.Name, payload.Title, payload.Message)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
