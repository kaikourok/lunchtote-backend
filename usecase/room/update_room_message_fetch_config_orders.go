package room

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *RoomUsecase) UpdateRoomMessageFetchConfigOrders(characterId int, orders *[]model.RoomMessageFetchConfigOrder) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.UpdateRoomMessageFetchConfigOrders(characterId, orders)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
