package diary

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (s *DiaryUsecase) RetrieveDiaryEditData(characterId int) (*model.DiaryEditData, error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	data, err := repository.RetrieveDiaryEditData(characterId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return data, nil
}
