package character

func (s *CharacterUsecase) RetrieveInitialData(id int) (existsUnreadNotification, existsUnreadMail bool, err error) {
	repository := s.registry.GetRepository()
	return repository.RetrieveInitialData(id)
}
