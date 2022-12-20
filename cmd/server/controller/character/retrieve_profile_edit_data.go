package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RetrieveProfileEditData(c *gin.Context) {
	session := sessions.Default(c)

	character, err := u.usecase.RetrieveProfileEditData(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, character)
}
