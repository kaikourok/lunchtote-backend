package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) CreateList(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Name string `json:"name"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	listId, err := u.usecase.CreateList(session.Get("cid").(int), payload.Name)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"listId": listId,
	})
}
