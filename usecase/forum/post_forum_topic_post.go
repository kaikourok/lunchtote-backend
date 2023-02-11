package forum

import (
	"errors"

	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/library/secure"
)

func (s *ForumUsecase) PostForumTopicPost(characterId *int, isAdministrator *bool, ip *string, topicId int, post *model.ForumTopicPostSendData) (postId int, err error) {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	if post.PostType == "ADMINISTRATOR" {
		if isAdministrator == nil || !*isAdministrator {
			return 0, errors.New("管理者ではありません")
		}
	}

	savePost := *post
	if savePost.EditPassword != nil {
		hashed, err := secure.HashPassword(*savePost.EditPassword, config.GetInt("secure.bcrypt-cost"))
		if err != nil {
			return 0, err
		}
		savePost.EditPassword = &hashed
	}

	postId, err = repository.PostForumTopicPost(characterId, s.generateIdentifier(ip), topicId, &savePost)
	if err != nil {
		logger.Error(err)
		return
	}

	return postId, nil
}
