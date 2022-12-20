package general

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *GeneralController) Connections(c *gin.Context) {
	connections, err := u.usecase.Connections()

	if err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"connections": connections,
		})
	}
}
