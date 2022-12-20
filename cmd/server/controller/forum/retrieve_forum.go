package forum

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (u *ForumController) RetrieveForum(c *gin.Context) {
	forumId, err := strconv.Atoi(c.Param("forum"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	forum, err := u.usecase.RetrieveForum(forumId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, forum)
}
