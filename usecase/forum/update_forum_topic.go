package forum

import (
	"errors"

	"github.com/kaikourok/lunchtote-backend/entity/model"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *ForumUsecase) UpdateForumTopic(characterId *int, isAdministrator *bool, topicId int, topic *model.ForumTopicEditData) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	masterCharacter, editPassword, postType, err := repository.RetrieveForumTopicEditCredentials(topicId)
	if err != nil {
		logger.Error(err)
		return err
	}

	if postType != topic.PostType {
		return errors.New("投稿タイプが一致していません")
	}

	switch topic.PostType {
	case "ANONYMOUS":
		if topic.EditPassword == nil {
			return usecaseErrors.ErrValidate
		}

		err = bcrypt.CompareHashAndPassword([]byte(*editPassword), []byte(*topic.EditPassword))
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

	err = repository.UpdateForumTopic(characterId, topicId, topic)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
