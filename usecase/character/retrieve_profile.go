package character

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) RetrieveProfile(characterId *int, targetId int) (*model.Profile, error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(targetId, validation.Min(1))
	if err != nil {
		return nil, usecaseErrors.ErrValidate
	}

	profile, err := repository.RetrieveProfile(characterId, targetId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if profile.IsBlocked {
		err := errors.New("ブロックされているため情報を表示できません")
		return nil, err
	}

	return profile, nil
}
