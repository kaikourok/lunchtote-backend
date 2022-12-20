package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RetrieveLayeringGroupOverviews(c *gin.Context) {
	session := sessions.Default(c)

	overviews, err := u.usecase.RetrieveLayeringGroupOverviews(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"layeringGroups": overviews,
	})
}
