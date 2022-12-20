package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) DeleteList(userId int, listId int) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(listId, validation.Min(1))
	if err != nil {
		return errors.ErrValidate
	}

	err = repository.DeleteList(userId, listId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
