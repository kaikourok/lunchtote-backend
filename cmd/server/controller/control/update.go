package control

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *ControlController) Update(c *gin.Context) {

	err := u.usecase.Update()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
