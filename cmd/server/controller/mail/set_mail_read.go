package mail

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *MailController) SetMailRead(c *gin.Context) {
	session := sessions.Default(c)

	mailId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	existsUnreadMail, err := u.usecase.SetMailRead(session.Get("cid").(int), mailId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"existsUnreadMail": existsUnreadMail,
	})
}
