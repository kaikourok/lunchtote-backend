package character

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (u *CharacterController) UpdateLayerGroupOrders(c *gin.Context) {
	session := sessions.Default(c)

	layeringGroup, err := strconv.Atoi(c.Param("layering_group"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var payload struct {
		Orders *[]model.CharacterIconLayerGroupOrder `json:"orders"`
	}
	err = c.BindJSON(&payload)
	if err != nil {
		return
	}

	err = u.usecase.UpdateLayerGroupOrders(session.Get("cid").(int), layeringGroup, payload.Orders)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
