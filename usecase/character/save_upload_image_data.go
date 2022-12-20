package character

import (
	"bytes"

	"github.com/kaikourok/lunchtote-backend/entity/service"
)

func (s *CharacterUsecase) SaveUploadImageData(characterId int, imageBuffers []*bytes.Buffer, imageType service.ImageTypeId) (paths *[]string, err error) {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	paths, err = repository.SaveUploadImageData(characterId, imageBuffers, imageType, config.GetString("general.upload-directory"))
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return paths, nil
}
