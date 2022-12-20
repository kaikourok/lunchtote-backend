package forum

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *ForumUsecase) RetrieveForumOverviews() (forumGroups *[]model.ForumGroup, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	forumGroups, err = repository.RetrieveForumOverviews()
	if err != nil {
		logger.Error(err)
		return
	}

	return forumGroups, nil
}
