package mail

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *MailController) RetrieveMails(c *gin.Context) {
	session := sessions.Default(c)

	start := math.MaxInt32
	if c.Query("start") != "" {
		parsedStart, err := strconv.Atoi(c.Query("start"))
		if err == nil && 0 <= parsedStart {
			start = parsedStart
		}
	}

	unreadOnly := c.Query("unread") == "t"

	mails, isContinue, err := u.usecase.RetrieveMails(session.Get("cid").(int), start, unreadOnly)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mails":      mails,
		"isContinue": isContinue,
	})
}
