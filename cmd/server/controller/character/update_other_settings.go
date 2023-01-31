package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *CharacterController) UpdateOtherSettings(c *gin.Context) {
	session := sessions.Default(c)

	var payload model.CharacterOtherSettings
	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.UpdateOtherSettings(session.Get("cid").(int), &payload)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
