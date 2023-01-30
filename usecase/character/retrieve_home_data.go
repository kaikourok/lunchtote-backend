package character

import (
	"math"

	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *CharacterUsecase) RetrieveHomeData(characterId int) (home *model.HomeData, announcements *[]model.AnnouncementOverview, err error) {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	home, err = repository.RetrieveHomeData(characterId)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	announcements, _, err = repository.RetrieveAnnouncementOverviews(math.MaxInt32, config.GetInt("general.home-announcements-max"))
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	return home, announcements, nil
}
