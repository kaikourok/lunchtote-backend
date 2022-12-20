package diary

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *DiaryController) ClearReservedDiary(c *gin.Context) {
	session := sessions.Default(c)

	err := u.usecase.ClearReservedDiary(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
