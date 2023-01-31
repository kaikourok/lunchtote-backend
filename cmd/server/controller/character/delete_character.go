package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) DeleteCharacter(c *gin.Context) {
	session := sessions.Default(c)

	err := u.usecase.DeleteCharacter(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	session.Clear()
	session.Options(u.getSessionDefaultConfig())
	session.Save()

	c.Status(http.StatusOK)
}
