package forum

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (u *ForumController) RetrieveForumForcedPostType(c *gin.Context) {
	forum, err := strconv.Atoi(c.Param("forum"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	forcedPostType, err := u.usecase.RetrieveForumForcedPostType(forum)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"forcedPostType": forcedPostType,
	})
}
