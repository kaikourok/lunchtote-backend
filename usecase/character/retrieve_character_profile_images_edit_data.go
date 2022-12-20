package character

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *CharacterUsecase) RetrieveCharacterProfileImagesEditData(characterId int) (*[]model.ProfileImage, error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	images, err := repository.RetrieveCharacterProfileImages(characterId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return images, nil
}
