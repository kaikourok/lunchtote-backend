package character

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *CharacterController) UpdatePasswordByResetCode(c *gin.Context) {
	var payload struct {
		Id       int    `json:"id"`
		Code     string `json:"code"`
		Password string `json:"password"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.UpdatePasswordByResetCode(payload.Id, payload.Code, payload.Password)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			c.Status(http.StatusForbidden)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}
