package character

import (
	"math"

	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *CharacterUsecase) RetrieveHomeData(characterId int) (home *model.HomeData, announcements *[]model.AnnouncementOverview, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	home, err = repository.RetrieveHomeData(characterId)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	announcements, _, err = repository.RetrieveAnnouncementOverviews(math.MaxInt32, 5)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	return home, announcements, nil
}
