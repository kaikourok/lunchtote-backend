package forum

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *ForumUsecase) RetrieveForumTopicOverviews(forumId, page int) (topics *[]model.ForumTopicOverview, pages int, err error) {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	forumTopicsPerPage := config.GetInt("general.forum-topics-per-page")

	topics, topicCounts, err := repository.RetrieveForumTopicOverviews(forumId, (page-1)*forumTopicsPerPage, forumTopicsPerPage)
	if err != nil {
		logger.Error(err)
		return
	}

	pages = topicCounts / forumTopicsPerPage
	if topicCounts%forumTopicsPerPage != 0 {
		pages++
	}

	return topics, pages, nil
}
