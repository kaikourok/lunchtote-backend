package general

import (
	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/library/database"
)

type GeneralRepository struct {
	*sqlx.DB
}

func (r *GeneralRepository) ExecTx(f func(tx *sqlx.Tx) error) error {
	return database.ExecTx(r, f)
}

func NewGeneralRepository(db *sqlx.DB) *GeneralRepository {
	return &GeneralRepository{db}
}
