package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RemoveCharacterFromList(c *gin.Context) {
	session := sessions.Default(c)

	list, err := strconv.Atoi(c.Param("list"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var payload struct {
		Character *int `json:"character"`
	}

	err = c.BindJSON(&payload)
	if err != nil {
		return
	}

	err = u.usecase.RemoveCharacterFromList(session.Get("cid").(int), *payload.Character, list)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
