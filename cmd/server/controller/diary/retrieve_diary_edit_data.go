package diary

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *DiaryController) RetrieveDiaryEditData(c *gin.Context) {
	session := sessions.Default(c)

	diary, err := u.usecase.RetrieveDiaryEditData(session.Get("cid").(int))
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, diary)
}
