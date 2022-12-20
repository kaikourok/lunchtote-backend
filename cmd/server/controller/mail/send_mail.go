package mail

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *MailController) SendMail(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Target  int    `json:"target"`
		Title   string `json:"title"`
		Message string `json:"message"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.SendMail(session.Get("cid").(int), payload.Target, payload.Title, payload.Message)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
