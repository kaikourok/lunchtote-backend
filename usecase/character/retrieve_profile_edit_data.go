package character

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *CharacterUsecase) RetrieveProfileEditData(characterId int) (*model.ProfileEditData, error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	data, err := repository.RetrieveProfileEditData(characterId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return data, nil
}
