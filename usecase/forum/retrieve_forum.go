package forum

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *ForumUsecase) RetrieveForum(forumId int) (forum *model.Forum, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	forum, err = repository.RetrieveForum(forumId)
	if err != nil {
		logger.Error(err)
		return
	}

	return forum, nil
}
