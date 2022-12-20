package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RetrieveCharacterList(c *gin.Context) {
	session := sessions.Default(c)

	page := 0
	if c.Query("page") != "" {
		parsedPage, err := strconv.Atoi(c.Query("page"))
		if err == nil && 0 <= parsedPage {
			page = parsedPage
		}
	}

	var characterId *int
	if session.Get("cid") != nil {
		v := session.Get("cid").(int)
		characterId = &v
	}

	list, maxId, err := u.usecase.RetrieveCharacterList(characterId, page)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"characters": list,
		"maxId":      maxId,
	})
}
