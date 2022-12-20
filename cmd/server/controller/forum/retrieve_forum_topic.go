package forum

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *ForumController) RetrieveForumTopic(c *gin.Context) {
	session := sessions.Default(c)

	topicId, err := strconv.Atoi(c.Param("topic"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	topic, err := u.usecase.RetrieveForumTopic(topicId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	var characterId *int
	if session.Get("cid") != nil {
		v := session.Get("cid").(int)
		characterId = &v
	}

	posts, err := u.usecase.RetrieveForumTopicPosts(topicId, characterId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"topic": topic,
		"posts": posts,
	})
}
