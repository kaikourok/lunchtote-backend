package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *CharacterController) UpdateIcons(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Icons      []model.Icon `json:"icons"`
		InsertOnly bool         `json:"insertOnly,omitempty"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.UpdateIcons(session.Get("cid").(int), &payload.Icons, payload.InsertOnly)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
