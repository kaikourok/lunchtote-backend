package diary

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *DiaryController) RetrieveDiaryOverviews(c *gin.Context) {
	session := sessions.Default(c)

	var nth *int
	if c.Query("nth") != "" {
		parsedNth, err := strconv.Atoi(c.Query("nth"))
		if err == nil && 0 < parsedNth {
			nth = &parsedNth
		}
	}

	diaries, currentNth, lastNth, err := u.usecase.RetrieveDiaryOverviews(session.Get("cid").(int), nth)
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"diaries":    diaries,
		"currentNth": currentNth,
		"lastNth":    lastNth,
	})
}
