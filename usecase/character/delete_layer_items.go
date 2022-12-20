package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) DeleteLayerItems(characterId int, itemIds *[]int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	if itemIds == nil {
		return errors.ErrValidate
	}

	err := validation.Validate(*itemIds, validation.Required, validation.Each(
		validation.Min(1),
	))
	if err != nil {
		return errors.ErrValidate
	}

	err = repository.DeleteLayerItems(characterId, itemIds)
	if err != nil {
		logger.Error(err)
	}

	return err
}
