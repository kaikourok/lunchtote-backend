package diary

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *DiaryController) RetrieveDiaryPreview(c *gin.Context) {
	session := sessions.Default(c)

	diary, err := u.usecase.RetrieveDiaryPreview(session.Get("cid").(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if diary == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, diary)
}
