package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RetrieveNotifications(c *gin.Context) {
	session := sessions.Default(c)

	start := 0
	if c.Query("start") != "" {
		parsedStart, err := strconv.Atoi(c.Query("start"))
		if err == nil && 0 <= parsedStart {
			start = parsedStart
		}
	}

	notifications, isContinue, err := u.usecase.RetrieveNotifications(session.Get("cid").(int), start)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"isContinue":    isContinue,
	})
}
