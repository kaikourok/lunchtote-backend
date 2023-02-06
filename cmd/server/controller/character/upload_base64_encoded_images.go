package character

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/service"
)

func (u *CharacterController) UploadBase64EncordedImages(c *gin.Context) {
	session := sessions.Default(c)

	var imageType service.ImageTypeId
	switch c.Query("type") {
	case "icon":
		imageType = service.ImageType.Icon
	case "icon-fragment":
		imageType = service.ImageType.IconFragment
	case "profile-image":
		imageType = service.ImageType.ProfileImage
	case "list-image":
		imageType = service.ImageType.ListImage
	default:
		imageType = service.ImageType.General
	}

	var payload struct {
		Images []string `json:"images"`
	}

	err := c.BindJSON(&payload)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	imageBuffers := make([]*bytes.Buffer, 0, len(payload.Images))
	for i := range payload.Images {
		_, data, ok := strings.Cut(payload.Images[i], ",")
		if !ok {
			c.Status(http.StatusBadRequest)
			return
		}

		imageBytes, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		imageBuffers = append(imageBuffers, bytes.NewBuffer(imageBytes))
	}

	paths, err := u.usecase.SaveUploadImageData(session.Get("cid").(int), imageBuffers, imageType)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"paths": paths,
	})
}
