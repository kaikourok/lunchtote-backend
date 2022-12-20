package forum

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *ForumController) UpdateForumTopic(c *gin.Context) {
	session := sessions.Default(c)

	topic, err := strconv.Atoi(c.Param("topic"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var payload model.ForumTopicEditData
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

	err = u.usecase.UpdateForumTopic(characterId, isAdministrator, topic, &payload)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
