package room

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *RoomUsecase) AddRoomMessageFetchConfig(characterId int, config *model.RoomMessageFetchConfig) (configId int, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	configId, err = repository.AddRoomMessageFetchConfig(characterId, config)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return configId, nil
}
