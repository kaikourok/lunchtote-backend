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

	nameChannnel := make(chan nameResultStruct)
	iconsChannel := make(chan iconsResultStruct)

	go func() {
		name, err := repository.RetrieveCharacterNickname(characterId)
		nameChannnel <- nameResultStruct{name, err}
	}()

	go func() {
		icons, err := repository.RetrieveCharacterIcons(characterId)
		iconsChannel <- iconsResultStruct{icons, err}
	}()

	nameResult := <-nameChannnel
	iconsResult := <-iconsChannel

	if nameResult.err != nil || iconsResult.err != nil {
		if nameResult.err != nil {
			logger.Error(nameResult.err)
		}
		if iconsResult.err != nil {
			logger.Error(iconsResult.err)
		}
		return nil, errors.New("データの取得中にエラーが発生しました")
	}

	var result model.RoomMessageEditRequiredData
	result.Character.Name = nameResult.name
	result.Icons = *iconsResult.icons

	return &result, nil
}
