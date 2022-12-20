package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/library/secure"
)

func (u *CharacterController) SignIn(c *gin.Context) {
	config := u.registry.GetConfig()
	session := sessions.Default(c)

	var payload struct {
		Key      string `json:"key"`
		Password string `json:"password"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	id, notificationToken, administrator, err := u.usecase.SignIn(payload.Key, payload.Password)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}

	csrfToken := secure.GenerateSecureRandomHex(config.GetInt("secure.token-length"))

	session.Set("cid", id)
	session.Set("csrf-token", csrfToken)
	session.Set("notification-token", notificationToken)
	session.Set("administrator", administrator)
	session.Options(u.getSessionDefaultConfig())
	session.Save()

	c.Status(http.StatusOK)
}
