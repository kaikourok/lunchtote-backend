package room

import (
	"errors"

	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *RoomUsecase) RetrieveRoomMessageEditRequiredData(characterId int) (data *model.RoomMessageEditRequiredData, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	type nameResultStruct struct {
		name string
		err  error
	}

	type iconsResultStruct struct {
		icons *[]model.Icon
		err   error
	}

	type listsResultStruct struct {
		lists *[]model.ListOverview
		err   error
	}

	type fetchConfigResultStruct struct {
		config *[]model.RoomMessageFetchConfig
		err    error
	}

	nameChannnel := make(chan nameResultStruct)
	iconsChannel := make(chan iconsResultStruct)
	listsChannel := make(chan listsResultStruct)
	fetchConfigChannel := make(chan fetchConfigResultStruct)

	go func() {
		name, err := repository.RetrieveCharacterNickname(characterId)
		nameChannnel <- nameResultStruct{name, err}
	}()

	go func() {
		icons, err := repository.RetrieveCharacterIcons(characterId)
		iconsChannel <- iconsResultStruct{icons, err}
	}()

	go func() {
		lists, err := repository.RetrieveLists(characterId)
		listsChannel <- listsResultStruct{lists, err}
	}()

	go func() {
		config, err := repository.RetrieveRoomMessageFetchConfig(characterId)
		fetchConfigChannel <- fetchConfigResultStruct{config, err}
	}()

	nameResult := <-nameChannnel
	iconsResult := <-iconsChannel
	listsResult := <-listsChannel
	fetchConfigResult := <-fetchConfigChannel

	if nameResult.err != nil || iconsResult.err != nil || listsResult.err != nil || fetchConfigResult.err != nil {
		if nameResult.err != nil {
			logger.Error(nameResult.err)
		}
		if iconsResult.err != nil {
			logger.Error(iconsResult.err)
		}
		if listsResult.err != nil {
			logger.Error(listsResult.err)
		}
		if fetchConfigResult.err != nil {
			logger.Error(fetchConfigResult.err)
		}
		return nil, errors.New("データの取得中にエラーが発生しました")
	}

	var result model.RoomMessageEditRequiredData
	result.Character.Name = nameResult.name
	result.Icons = *iconsResult.icons
	result.Lists = *listsResult.lists
	result.FetchConfigs = *fetchConfigResult.config

	return &result, nil
}
