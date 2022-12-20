package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) SignOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(u.getSessionDefaultConfig())
	session.Save()
	c.Status(http.StatusOK)
}
