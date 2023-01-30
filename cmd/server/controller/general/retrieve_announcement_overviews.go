package general

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *GeneralController) RetrieveAnnouncementOverviews(c *gin.Context) {
	announcements, err := u.usecase.RetrieveAnnouncementOverviews()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, announcements)
}
