package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) DeleteUploadedImages(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Targets []int `json:"targets"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.DeleteUploadedImages(session.Get("cid").(int), &payload.Targets)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
