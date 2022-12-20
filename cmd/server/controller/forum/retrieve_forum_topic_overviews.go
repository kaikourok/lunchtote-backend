package forum

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (u *ForumController) RetrieveForumTopicOverviews(c *gin.Context) {
	forum, err := strconv.Atoi(c.Param("forum"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}

	topics, pages, err := u.usecase.RetrieveForumTopicOverviews(forum, page)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"topics": topics,
		"pages":  pages,
	})
}
