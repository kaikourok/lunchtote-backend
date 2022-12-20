package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) RetrieveCharacterProfileImagesEditData(c *gin.Context) {
	session := sessions.Default(c)

	profileImages, err := u.usecase.RetrieveCharacterProfileImagesEditData(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, profileImages)
}
