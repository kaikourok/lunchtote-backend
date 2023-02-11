package diary

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/entity/service"
)

func (s *DiaryUsecase) RetrieveDiaryPreview(characterId int) (*model.Diary, error) {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	diary, err := repository.RetrieveDiaryPreview(characterId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if diary == nil {
		return nil, nil
	}

	diary.Diary = service.StylizeTextEntry(diary.Diary, config.GetString("general.upload-directory"))

	return diary, nil
}
