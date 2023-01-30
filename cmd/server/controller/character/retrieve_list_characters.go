package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RetrieveList(c *gin.Context) {
	session := sessions.Default(c)

	listId, err := strconv.Atoi(c.Param("list"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	name, characters, err := u.usecase.RetrieveList(session.Get("cid").(int), listId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":       name,
		"characters": characters,
	})
}
