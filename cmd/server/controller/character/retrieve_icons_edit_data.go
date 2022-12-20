package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RetrieveCharacterIconsEditData(c *gin.Context) {
	session := sessions.Default(c)

	icons, err := u.usecase.RetrieveCharacterIconsEditData(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, icons)
}
