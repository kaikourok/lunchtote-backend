package forum

import (
	"errors"

	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *ForumUsecase) DeleteForumTopicPost(characterId *int, isAdministrator *bool, postId int, editPassword *string) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	masterCharacter, savedPassword, postType, err := repository.RetrieveForumTopicEditCredentials(postId)
	if err != nil {
		logger.Error(err)
		return err
	}

	switch postType {
	case "ANONYMOUS":
		if editPassword == nil {
			return usecaseErrors.ErrValidate
		}

		err = bcrypt.CompareHashAndPassword([]byte(*savedPassword), []byte(*editPassword))
		if err != nil {
			return usecaseErrors.ErrPermission
		}
	case "SIGNED_IN":
		if characterId == nil || *characterId != *masterCharacter {
			return usecaseErrors.ErrPermission
		}
	case "ADMINISTRATOR":
		if isAdministrator == nil || !*isAdministrator {
			return usecaseErrors.ErrPermission
		}
	default:
		return errors.New("保存されているデータが正しくありません")
	}

	err = repository.DeleteForumTopicPost(characterId, postId)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
