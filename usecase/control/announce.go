package control

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (s *ControlUsecase) Announce(announce *model.AnnouncementData) error {
	repository := s.registry.GetRepository()
	logger := s.registry.GetLogger()

	err := repository.Announce(announce)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
