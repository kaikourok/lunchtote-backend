package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RenameList(c *gin.Context) {
	session := sessions.Default(c)

	listId, err := strconv.Atoi(c.Param("list"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var payload struct {
		Name string `json:"name"`
	}
	err = c.BindJSON(&payload)
	if err != nil {
		return
	}

	err = u.usecase.RenameList(session.Get("cid").(int), listId, payload.Name)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
