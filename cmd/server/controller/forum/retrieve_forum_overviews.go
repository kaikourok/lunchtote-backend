package forum

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *ForumController) RetrieveForumOverviews(c *gin.Context) {
	overview, err := u.usecase.RetrieveForumOverviews()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, overview)
}
