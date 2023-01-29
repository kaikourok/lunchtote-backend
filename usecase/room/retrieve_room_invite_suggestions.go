package room

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) RetrieveRoomInviteSuggestions(characterId int, searchText string, roomId int) (suggestions *model.CharacterSuggestions, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(roomId, validation.Min(1))
	if err != nil {
		return nil, usecaseErrors.ErrValidate
	}

	permissions, _, banned, err := repository.RetrieveRoomOwnPermissions(characterId, roomId)
	if err != nil {
		return nil, err
	}
	if banned || !permissions.Invite {
		return nil, usecaseErrors.ErrPermission
	}

	suggestionsData, err := repository.RetrieveRoomInviteSuggestions(characterId, searchText, roomId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return suggestionsData.ToDomain(), nil
}
