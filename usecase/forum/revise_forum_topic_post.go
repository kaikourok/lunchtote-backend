package forum

import (
	"errors"

	"github.com/kaikourok/lunchtote-backend/entity/model"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *ForumUsecase) ReviseForumTopicPost(characterId *int, isAdministrator *bool, postId int, post *model.ForumTopicPostReviseData) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	masterCharacter, editPassword, postType, err := repository.RetrieveForumTopicEditCredentials(postId)
	if err != nil {
		logger.Error(err)
		return err
	}

	if postType != post.PostType {
		return errors.New("投稿タイプが一致していません")
	}

	switch post.PostType {
	case "ANONYMOUS":
		if post.EditPassword == nil {
			return usecaseErrors.ErrValidate
		}

		err = bcrypt.CompareHashAndPassword([]byte(*editPassword), []byte(*post.EditPassword))
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

	err = repository.ReviseForumTopicPost(characterId, postId, post)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
