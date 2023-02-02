package control

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *ControlController) UpdateAnnouncement(c *gin.Context) {
	announcementId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var payload model.AnnouncementEditDataUpdate
	err = c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.UpdateAnnouncement(announcementId, &payload)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
