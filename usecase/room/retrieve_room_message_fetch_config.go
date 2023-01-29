package room

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *RoomUsecase) RetrieveRoomMessageFetchConfig(characterId int) (configs *[]model.RoomMessageFetchConfigWithId, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	configs, err = repository.RetrieveRoomMessageFetchConfig(characterId)
	if err != nil {
		logger.Error(err)
	}

	return
}
