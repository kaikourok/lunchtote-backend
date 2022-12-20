package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *CharacterController) DeleteLayerGroup(c *gin.Context) {
	session := sessions.Default(c)

	layerGroup, err := strconv.Atoi(c.Param("layer_group"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = u.usecase.DeleteLayerGroup(session.Get("cid").(int), layerGroup)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
