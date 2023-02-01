package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) UpdateNotificationChecked(c *gin.Context) {
	session := sessions.Default(c)

	err := u.usecase.UpdateNotificationChecked(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
