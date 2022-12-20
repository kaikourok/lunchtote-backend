package diary

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *DiaryController) RetrieveDiary(c *gin.Context) {
	session := sessions.Default(c)

	nth, err := strconv.Atoi(c.Param("nth"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var characterId *int
	if session.Get("cid") != nil {
		cid := session.Get("cid").(int)
		characterId = &cid
	}

	targetId, err := strconv.Atoi(c.Param("character"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	diary, err := u.usecase.RetrieveDiary(characterId, targetId, nth)
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, diary)
}
