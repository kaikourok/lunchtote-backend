package repository

import "github.com/kaikourok/lunchtote-backend/entity/model"

type diaryRepository interface {
	// 日記取得関連
	RetrieveLatestDiaryNth() (nth int, err error)
	RetrieveDiary(characterId *int, targetId, nth int) (*model.Diary, error)
	RetrieveDiaryPreview(characterId int) (*model.Diary, error)
	RetrieveDiaryOverviews(characterId, nth int) (diaries *[]model.DiaryOverview, err error)

	// 投稿関連
	RetrieveDiaryEditData(characterId int) (*model.DiaryEditData, error)
	ReservePublishDiary(characterId int, title, diary string) error
	ClearReservedDiary(characterId int) error
}
