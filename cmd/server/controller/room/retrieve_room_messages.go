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

	if param := c.Query("list"); param != "" {
		if list, err := strconv.Atoi(param); err == nil {
			options = append(options, room.RetrieveRoomMessagesOptionListId(&list))
		}
	}

	if param := c.Query("character"); param != "" {
		if character, err := strconv.Atoi(param); err == nil {
			options = append(options, room.RetrieveRoomMessagesOptionTargetCharacterId(&character))
		}
	}

	if param := c.Query("root"); param != "" {
		if referRoot, err := strconv.Atoi(param); err == nil {
			options = append(options, room.RetrieveRoomMessagesOptionReferRoot(&referRoot))
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

	if param := c.Query("base"); param != "" {
		if base, err := strconv.Atoi(param); err == nil {
			options = append(options, room.RetrieveRoomMessagesOptionBasePoint(base))
		}
	}

	messages, isContinueFollowing, isContinuePrevious, err := u.usecase.RetrieveRoomMessages(session.Get("cid").(int), options...)
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
