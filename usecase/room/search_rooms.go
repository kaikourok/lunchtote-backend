package room

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) SearchRooms(characterId int, options *model.RoomSearchOptions) (rooms []model.RoomListItem, isContinue bool, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.ValidateStruct(options,
		validation.Field(&options.Order, validation.In("latest-post", "id", "posts-per-day")),
		validation.Field(&options.Sort, validation.In("asc", "desc")),
		validation.Field(&options.Participant, validation.In(nil, "own", "follow")),
	)
	if err != nil {
		return nil, false, usecaseErrors.ErrValidate
	}

	rooms, isContinue, err = repository.SearchRooms(characterId, options)
	if err != nil {
		logger.Error(err)
		return nil, false, err
	}

	return
}
