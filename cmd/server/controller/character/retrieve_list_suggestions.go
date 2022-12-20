package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RetrieveListSuggestions(c *gin.Context) {
	session := sessions.Default(c)

	list, err := strconv.Atoi(c.Param("list"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var payload struct {
		Text string `json:"text"`
	}

	err = c.BindJSON(&payload)
	if err != nil {
		return
	}

	results, err := u.usecase.RetrieveListSuggestions(session.Get("cid").(int), payload.Text, list)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, results)
}
