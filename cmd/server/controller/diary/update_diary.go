package diary

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *DiaryController) ReservePublishDiary(c *gin.Context) {
	session := sessions.Default(c)

	var payload struct {
		Title string `json:"title"`
		Diary string `json:"diary"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		return
	}

	err = u.usecase.ReservePublishDiary(session.Get("cid").(int), payload.Title, payload.Diary)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
