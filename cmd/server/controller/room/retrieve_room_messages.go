package room

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/usecase/room"
)

func (u *RoomController) RetrieveRoomMessages(c *gin.Context) {
	session := sessions.Default(c)

	/*roomId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}*/

	options := make([]room.RetrieveRoomMessagesOption, 0, 16)

	if param := c.Query("type"); param != "" {
		options = append(options, room.RetrieveRoomMessagesOptionRangeType(param))
	}

	if param := c.Query("category"); param != "" {
		options = append(options, room.RetrieveRoomMessagesOptionCategory(param))
	}

	if param := c.Query("room"); param != "" {
		if roomId, err := strconv.Atoi(param); err == nil {
			options = append(options, room.RetrieveRoomMessagesOptionRoomId(&roomId))
		}
	}

	if param := c.Query("relates"); param != "" {
		options = append(options, room.RetrieveRoomMessagesOptionRelateFilter(param == "true"))
	}

	if param := c.Query("children"); param != "" {
		options = append(options, room.RetrieveRoomMessagesOptionChildren(param == "true"))
	}

	if param := c.Query("search"); param != "" {
		options = append(options, room.RetrieveRoomMessagesOptionSearch(&param))
	}

	messages, isContinuePrevious, isContinueFollowing, err := u.usecase.RetrieveRoomMessages(session.Get("cid").(int), options...)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages":            messages,
		"isContinueFollowing": isContinueFollowing,
		"isContinuePrevious":  isContinuePrevious,
	})
}
