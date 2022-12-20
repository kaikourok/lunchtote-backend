package character

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/entity/service"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) UpdateLayerItems(characterId, groupId int, items *[]model.CharacterIconLayerItemEditData) (result *[]model.CharacterIconLayerItem, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(groupId, validation.Min(1))
	if err != nil {
		return nil, usecaseErrors.ErrValidate
	}

	if items == nil {
		return nil, usecaseErrors.ErrValidate
	}

	for _, item := range *items {
		character, err := service.ParseFilePath(item.Path)
		if err != nil {
			return nil, usecaseErrors.ErrValidate
		}
		if character != characterId {
			return nil, errors.New("自身のものではない画像を参照しています")
		}
	}

	result, err = repository.UpdateLayerItems(characterId, groupId, items)
	if err != nil {
		logger.Error(err)
	}

	return result, err
}
