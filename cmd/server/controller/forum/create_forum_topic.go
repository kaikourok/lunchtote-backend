package forum

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *ForumController) CreateForumTopic(c *gin.Context) {
	session := sessions.Default(c)

	forum, err := strconv.Atoi(c.Param("forum"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var payload model.ForumTopicCreateData
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

	var ip *string
	if payload.PostType == "ANONYMOUS" {
		ipEntity := c.ClientIP()
		ip = &ipEntity
	}

	topicId, err := u.usecase.CreateForumTopic(characterId, isAdministrator, ip, forum, &payload)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": topicId,
	})
}
