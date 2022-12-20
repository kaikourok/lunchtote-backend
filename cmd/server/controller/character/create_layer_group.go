package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) CreateLayerGroup(c *gin.Context) {
	session := sessions.Default(c)

	layeringGroup, err := strconv.Atoi(c.Param("layering_group"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var payload struct {
		Name string `json:"name"`
	}

	err = c.BindJSON(&payload)
	if err != nil {
		return
	}

	id, err := u.usecase.CreateLayerGroup(session.Get("cid").(int), layeringGroup, payload.Name)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
