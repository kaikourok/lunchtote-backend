package general

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (u *GeneralController) RetrieveAnnouncement(c *gin.Context) {
	announcementId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	announcement, prevGuide, nextGuide, err := u.usecase.RetrieveAnnouncement(announcementId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"announcement": announcement,
		"prevGuide":    prevGuide,
		"nextGuide":    nextGuide,
	})
}
