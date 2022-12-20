package forum

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *ForumController) ChangeForumReactionState(c *gin.Context) {
	session := sessions.Default(c)

	post, err := strconv.Atoi(c.Param("post"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var payload struct {
		Emoji string `json:"emoji"`
		State bool   `json:"state"`
	}
	err = c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if payload.State {
		err := u.usecase.ReactForumPost(session.Get("cid").(int), post, payload.Emoji)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
	} else {
		err := u.usecase.CancelReactForumPost(session.Get("cid").(int), post, payload.Emoji)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.Status(http.StatusOK)
}
