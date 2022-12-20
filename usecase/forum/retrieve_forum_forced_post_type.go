package forum

func (s *ForumUsecase) RetrieveForumForcedPostType(forumId int) (forcedPostType *string, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	forcedPostType, err = repository.RetrieveForumForcedPostType(forumId)
	if err != nil {
		logger.Error(err)
		return
	}

	return forcedPostType, nil
}
