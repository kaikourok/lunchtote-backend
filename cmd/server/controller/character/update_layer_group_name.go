package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) UpdateLayerGroupName(c *gin.Context) {
	session := sessions.Default(c)

	layerGroup, err := strconv.Atoi(c.Param("layer_group"))
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

	err = u.usecase.UpdateLayerGroupName(session.Get("cid").(int), layerGroup, payload.Name)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
