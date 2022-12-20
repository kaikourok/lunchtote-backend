package forum

import (
	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/library/database"
)

type ForumRepository struct {
	*sqlx.DB
}

func (r *ForumRepository) ExecTx(f func(tx *sqlx.Tx) error) error {
	return database.ExecTx(r, f)
}

func NewForumRepository(db *sqlx.DB) *ForumRepository {
	return &ForumRepository{db}
}
