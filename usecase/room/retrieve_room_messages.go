package room

import (
	"math"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

type RetrieveRoomMessagesOption func(*model.RoomMessageRetrieveSettings)

func RetrieveRoomMessagesOptionRangeType(rangeType string) RetrieveRoomMessagesOption {
	return func(args *model.RoomMessageRetrieveSettings) {
		args.RangeType = rangeType
	}
}

func RetrieveRoomMessagesOptionBasePoint(basePoint int) RetrieveRoomMessagesOption {
	return func(args *model.RoomMessageRetrieveSettings) {
		args.BasePoint = basePoint
	}
}

func RetrieveRoomMessagesOptionFetchNumber(fetchNumber int) RetrieveRoomMessagesOption {
	return func(args *model.RoomMessageRetrieveSettings) {
		args.FetchNumber = fetchNumber
	}
}

func RetrieveRoomMessagesOptionCategory(category string) RetrieveRoomMessagesOption {
	return func(args *model.RoomMessageRetrieveSettings) {
		args.Category = category
	}
}

func RetrieveRoomMessagesOptionRoomId(roomId *int) RetrieveRoomMessagesOption {
	return func(args *model.RoomMessageRetrieveSettings) {
		args.Room = roomId
	}
}

func RetrieveRoomMessagesOptionReferRoot(referRoot *int) RetrieveRoomMessagesOption {
	return func(args *model.RoomMessageRetrieveSettings) {
		args.ReferRoot = referRoot
	}
}

func RetrieveRoomMessagesOptionSearch(search *string) RetrieveRoomMessagesOption {
	return func(args *model.RoomMessageRetrieveSettings) {
		args.Search = search
	}
}

func RetrieveRoomMessagesOptionListId(listId *int) RetrieveRoomMessagesOption {
	return func(args *model.RoomMessageRetrieveSettings) {
		args.ListId = listId
	}
}

func RetrieveRoomMessagesOptionTargetCharacterId(targetCharacterId *int) RetrieveRoomMessagesOption {
	return func(args *model.RoomMessageRetrieveSettings) {
		args.TargetCharacterId = targetCharacterId
	}
}

func RetrieveRoomMessagesOptionRelateFilter(relateFilter bool) RetrieveRoomMessagesOption {
	return func(args *model.RoomMessageRetrieveSettings) {
		args.RelateFilter = relateFilter
	}
}

func RetrieveRoomMessagesOptionChildren(children bool) RetrieveRoomMessagesOption {
	return func(args *model.RoomMessageRetrieveSettings) {
		args.Children = children
	}
}

func (s *RoomUsecase) RetrieveRoomMessages(characterId int, options ...RetrieveRoomMessagesOption) (messages *[]model.RoomMessage, isContinuePrevious, isContinueFollowing *bool, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	option := &model.RoomMessageRetrieveSettings{
		RangeType:         "latest",
		BasePoint:         math.MaxInt32,
		FetchNumber:       100,
		Category:          "all",
		Room:              nil,
		ReferRoot:         nil,
		Search:            nil,
		ListId:            nil,
		TargetCharacterId: nil,
		RelateFilter:      false,
		Children:          true,
	}

	for _, optionSetter := range options {
		optionSetter(option)
	}

	err = validation.ValidateStruct(option,
		validation.Field(&option.RangeType, validation.In("latest", "initial", "previous", "following")),
		validation.Field(&option.Category, validation.In("all", "follow", "follow-other", "replied", "replied-other", "own", "conversation", "search", "list", "character", "character-replied")),
	)
	if err != nil {
		return nil, nil, nil, usecaseErrors.ErrValidate
	}

	switch option.Category {
	case "character":
	case "character-replied":
		if option.TargetCharacterId == nil {
			return nil, nil, nil, usecaseErrors.ErrValidate
		}
	case "conversation":
		if option.ReferRoot == nil {
			return nil, nil, nil, usecaseErrors.ErrValidate
		}
	case "search":
		if option.Search == nil || *option.Search == "" {
			return nil, nil, nil, usecaseErrors.ErrValidate
		}
	case "list":
		if option.ListId == nil {
			return nil, nil, nil, usecaseErrors.ErrValidate
		}
	}

	messages, isContinuePrevious, isContinueFollowing, err = repository.RetrieveRoomMessages(characterId, option)
	if err != nil {
		logger.Error(err)
		return nil, nil, nil, err
	}

	return
}
