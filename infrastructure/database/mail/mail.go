package mail

import (
	"github.com/jmoiron/sqlx"

	"github.com/kaikourok/lunchtote-backend/library/database"
)

type MailRepository struct {
	*sqlx.DB
}

func (r *MailRepository) ExecTx(f func(tx *sqlx.Tx) error) error {
	return database.ExecTx(r, f)
}

func NewMailRepository(db *sqlx.DB) *MailRepository {
	return &MailRepository{db}
}
