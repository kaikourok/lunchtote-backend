package forum

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *ForumUsecase) RetrieveForumTopicPosts(topicId int, characterId *int) (posts *[]model.ForumTopicPost, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	posts, err = repository.RetrieveForumTopicPosts(topicId, characterId)
	if err != nil {
		logger.Error(err)
		return
	}

	return posts, nil
}
