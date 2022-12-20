package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) RetrieveCredentialsByTwitter(twitterId string) (characterId int, notificationToken string, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err = validation.Validate(twitterId, validation.Required)
	if err != nil {
		return 0, "", errors.ErrValidate
	}

	characterId, notificationToken, err = repository.RetrieveCredentialsByTwitter(twitterId)
	if err != nil {
		logger.Error(err)
		return 0, "", err
	}

	return characterId, notificationToken, nil
}
