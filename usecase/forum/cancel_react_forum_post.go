package forum

func (s *ForumUsecase) CancelReactForumPost(characterId, postId int, emoji string) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := repository.CancelReactForumPost(characterId, postId, emoji)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
