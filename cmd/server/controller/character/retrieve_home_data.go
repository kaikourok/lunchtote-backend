package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RetrieveHomeData(c *gin.Context) {
	session := sessions.Default(c)

	home, announcements, err := u.usecase.RetrieveHomeData(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nickname":      home.Nickname,
		"announcements": announcements,
	})
}
