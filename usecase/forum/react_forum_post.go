package forum

func (s *ForumUsecase) ReactForumPost(characterId, postId int, emoji string) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.ReactForumPost(characterId, postId, emoji)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
