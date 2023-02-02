package control

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *ControlController) Announce(c *gin.Context) {
	var payload model.AnnouncementData
	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.Announce(&payload)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
