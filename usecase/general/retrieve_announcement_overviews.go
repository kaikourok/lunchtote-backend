package general

import (
	"math"

	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *GeneralUsecase) RetrieveAnnouncementOverviews() (announcements *[]model.AnnouncementOverview, err error) {
	repository := s.registry.GetRepository()
	logger := s.registry.GetLogger()

	// 本来ページネーション機能を持たせるべきだがひとまず全件取得
	announcements, _, err = repository.RetrieveAnnouncementOverviews(math.MaxInt32, 100000)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return
}
