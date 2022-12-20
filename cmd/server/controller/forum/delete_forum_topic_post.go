package forum

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *ForumController) DeleteForumTopicPost(c *gin.Context) {
	session := sessions.Default(c)

	post, err := strconv.Atoi(c.Param("post"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var payload struct {
		EditPassword *string `json:"editPassword"`
	}
	err = c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var characterId *int
	if session.Get("cid") != nil {
		v := session.Get("cid").(int)
		characterId = &v
	}

	var isAdministrator *bool
	if session.Get("administrator") != nil {
		v := session.Get("administrator").(bool)
		isAdministrator = &v
	}

	err = u.usecase.DeleteForumTopicPost(characterId, isAdministrator, post, payload.EditPassword)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
