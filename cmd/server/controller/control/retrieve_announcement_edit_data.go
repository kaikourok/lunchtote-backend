package control

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (u *ControlController) RetrieveAnnouncementEditData(c *gin.Context) {
	announcementId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	announcement, err := u.usecase.RetrieveAnnouncementEditData(announcementId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, announcement)
}
