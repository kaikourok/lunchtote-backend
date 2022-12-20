package general

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *GeneralController) Status(c *gin.Context) {
	c.JSON(http.StatusOK, &gin.H{
		"status": "working",
	})
}
