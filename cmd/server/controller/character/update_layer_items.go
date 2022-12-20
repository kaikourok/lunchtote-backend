package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *CharacterController) UpdateLayerItems(c *gin.Context) {
	session := sessions.Default(c)

	layerGroup, err := strconv.Atoi(c.Param("layer_group"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var payload struct {
		Images []model.CharacterIconLayerItemEditData `json:"images"`
	}

	err = c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	result, err := u.usecase.UpdateLayerItems(session.Get("cid").(int), layerGroup, &payload.Images)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, result)
}
