package forum

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *ForumUsecase) RetrieveForumTopic(topicId int) (topic *model.ForumTopic, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	topic, err = repository.RetrieveForumTopic(topicId)
	if err != nil {
		logger.Error(err)
		return
	}

	return topic, nil
}
