package character

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *CharacterUsecase) RetrieveUploadedImages(characterId int) (*[]model.UploadedImage, error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	images, err := repository.RetrieveUploadedImages(characterId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return images, nil
}
