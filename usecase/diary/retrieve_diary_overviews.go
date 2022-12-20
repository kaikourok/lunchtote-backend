package diary

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *DiaryUsecase) RetrieveDiaryOverviews(characterId int, nth *int) (diaries *[]model.DiaryOverview, currentNth, lastNth int, err error) {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	lastNth, err = repository.RetrieveLatestDiaryNth()
	if err != nil {
		logger.Error(err)
		return nil, 0, 0, err
	}

	currentNth = lastNth
	if nth != nil {
		currentNth = *nth
	}

	err = validation.Validate(currentNth, validation.Min(1))
	if err != nil {
		return nil, 0, 0, errors.ErrValidate
	}

	diaries, err = repository.RetrieveDiaryOverviews(characterId, currentNth)
	if err != nil {
		logger.Error(err)
		return nil, 0, 0, err
	}

	return diaries, currentNth, lastNth, nil
}
