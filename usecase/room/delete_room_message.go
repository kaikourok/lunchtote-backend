package room

import "github.com/kaikourok/lunchtote-backend/usecase/errors"

func (s *RoomUsecase) DeleteRoomMessage(characterId, messageId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	roomId, senderId, err := repository.RetrieveRoomMessageRelatedData(messageId)
	if err != nil {
		logger.Error(err)
		return err
	}

	if senderId != characterId {
		permissions, _, banned, err := repository.RetrieveRoomOwnPermissions(characterId, roomId)
		if err != nil {
			logger.Error(err)
			return err
		}
		if banned || !permissions.DeleteOtherMessage {
			return errors.ErrPermission
		}
	}

	err = repository.DeleteRoomMessage(messageId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
