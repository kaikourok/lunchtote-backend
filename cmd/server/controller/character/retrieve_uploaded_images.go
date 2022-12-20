package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RetrieveUploadedImages(c *gin.Context) {
	session := sessions.Default(c)

	images, err := u.usecase.RetrieveUploadedImages(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, images)
}
