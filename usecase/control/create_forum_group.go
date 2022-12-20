package control

func (s *ControlUsecase) CreateForumGroup(title string) (id int, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	id, err = repository.CreateForumGroup(title)
	if err != nil {
		logger.Error(err)
		return
	}

	return id, nil
}
