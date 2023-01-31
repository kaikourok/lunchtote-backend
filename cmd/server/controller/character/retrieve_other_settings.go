package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RetrieveOtherSettings(c *gin.Context) {
	session := sessions.Default(c)

	settings, err := u.usecase.RetrieveOtherSettings(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, settings)
}
