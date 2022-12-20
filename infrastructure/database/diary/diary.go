package diary

import (
	"github.com/jmoiron/sqlx"

	"github.com/kaikourok/lunchtote-backend/library/database"
)

type DiaryRepository struct {
	*sqlx.DB
}

func (db *DiaryRepository) ExecTx(f func(tx *sqlx.Tx) error) error {
	return database.ExecTx(db, f)
}

func NewDiaryRepository(db *sqlx.DB) *DiaryRepository {
	return &DiaryRepository{db}
}
