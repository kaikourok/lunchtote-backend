package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) DeleteLayerGroup(characterId, groupId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(groupId, validation.Min(1))
	if err != nil {
		return errors.ErrValidate
	}

	err = repository.DeleteLayerGroup(characterId, groupId)
	if err != nil {
		logger.Error(err)
	}

	return err
}
