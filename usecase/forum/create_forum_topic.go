package forum

import (
	"errors"

	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/library/secure"
)

func (s *ForumUsecase) CreateForumTopic(characterId *int, isAdministrator *bool, ip *string, forumId int, topic *model.ForumTopicCreateData) (topicId int, err error) {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	if topic.PostType == "ADMINISTRATOR" {
		if isAdministrator == nil || !*isAdministrator {
			return 0, errors.New("管理者ではありません")
		}
	}

	saveTopic := *topic
	if saveTopic.EditPassword != nil {
		hashed, err := secure.HashPassword(*saveTopic.EditPassword, config.GetInt("secure.bcrypt-cost"))
		if err != nil {
			return 0, err
		}
		saveTopic.EditPassword = &hashed
	}

	topicId, err = repository.CreateForumTopic(characterId, s.generateIdentifier(ip), forumId, &saveTopic)
	if err != nil {
		logger.Error(err)
		return
	}

	return topicId, nil
}
