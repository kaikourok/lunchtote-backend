package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RetrieveProfile(c *gin.Context) {
	session := sessions.Default(c)

	targetId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var characterId *int
	if session.Get("cid") == nil {
		*characterId = session.Get("cid").(int)
	}

	profile, err := u.usecase.RetrieveProfile(characterId, targetId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"character": profile,
	})
}
