package character

import (
	"bytes"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/entity/service"
)

func (u *CharacterController) UploadImages(c *gin.Context) {
	session := sessions.Default(c)

	var imageType service.ImageTypeId
	switch c.Query("type") {
	case "icon":
		imageType = service.ImageType.Icon
	case "icon-fragment":
		imageType = service.ImageType.IconFragment
	case "profile-image":
		imageType = service.ImageType.ProfileImage
	default:
		imageType = service.ImageType.General
	}

	form, _ := c.MultipartForm()
	images := form.File["images[]"]

	imageBuffers := make([]*bytes.Buffer, len(images))
	for i := range images {
		imageFile, err := images[i].Open()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		defer imageFile.Close()

		imageBuffer := &bytes.Buffer{}
		_, err = imageBuffer.ReadFrom(imageFile)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		imageBuffers[i] = imageBuffer
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
