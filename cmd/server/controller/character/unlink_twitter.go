package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) UnlinkTwitter(c *gin.Context) {
	session := sessions.Default(c)

	err := u.usecase.UnlinkTwitter(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
