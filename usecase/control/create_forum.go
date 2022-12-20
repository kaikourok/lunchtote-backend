package control

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *ControlUsecase) CreateForum(forumGroupId int, forum *model.ForumCreateData) (id int, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	id, err = repository.CreateForum(forumGroupId, forum)
	if err != nil {
		logger.Error(err)
		return
	}

	return id, nil
}
