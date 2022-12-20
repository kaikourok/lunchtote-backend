package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *CharacterController) UpdateProfileImages(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		ProfileImages []model.ProfileImage `json:"profileImages"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.UpdateProfileImages(session.Get("cid").(int), &payload.ProfileImages)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
