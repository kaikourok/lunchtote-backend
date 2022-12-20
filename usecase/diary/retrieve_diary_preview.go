package diary

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *DiaryUsecase) RetrieveDiaryPreview(characterId int) (*model.Diary, error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	diary, err := repository.RetrieveDiaryPreview(characterId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if diary == nil {
		return nil, nil
	}

	return diary, nil
}
