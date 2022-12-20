package character

func (s *CharacterUsecase) DeleteUploadedImages(characterId int, imageIds *[]int) error {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	err := repository.DeleteUploadedImages(characterId, imageIds, config.GetString("general.upload-directory"))
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
