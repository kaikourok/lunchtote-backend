package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RetrieveInitialData(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("cid").(int)

	existsUnreadNotification, existsUnreadMail, err := u.usecase.RetrieveInitialData(id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	response := struct {
		Id                       int    `json:"id"`
		CSRFToken                string `json:"csrfToken"`
		NotificationToken        string `json:"notificationToken"`
		ExistsUnreadNotification bool   `json:"existsUnreadNotification"`
		ExistsUnreadMail         bool   `json:"existsUnreadMail"`
		Administrator            bool   `json:"administrator,omitempty"`
	}{
		Id:                       id,
		CSRFToken:                session.Get("csrf-token").(string),
		NotificationToken:        session.Get("notification-token").(string),
		ExistsUnreadNotification: existsUnreadNotification,
		ExistsUnreadMail:         existsUnreadMail,
		Administrator:            session.Get("administrator").(bool),
	}

	c.JSON(http.StatusOK, response)
}
