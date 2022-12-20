package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RetrieveCharacterSuggestions(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Text       string `json:"text"`
		ExcludeOwn bool   `json:"excludeOwn"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	results, err := u.usecase.RetrieveCharacterSuggestions(session.Get("cid").(int), payload.Text, payload.ExcludeOwn)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, results)
}
