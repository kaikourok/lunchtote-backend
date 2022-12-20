package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) DeleteLayerItems(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Target []int `json:"target"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	err = u.usecase.DeleteLayerItems(session.Get("cid").(int), &payload.Target)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
