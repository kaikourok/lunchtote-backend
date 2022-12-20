package general

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *GeneralController) Inquiry(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Inquiry string `json:"inquiry"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	var characterId *int
	if session.Get("cid") != nil {
		c := session.Get("cid").(int)
		characterId = &c
	}

	err = u.usecase.Inquiry(characterId, payload.Inquiry)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.Status(http.StatusOK)
	}
}
