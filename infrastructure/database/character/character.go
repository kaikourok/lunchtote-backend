package character

import (
	"github.com/jmoiron/sqlx"

	"github.com/kaikourok/lunchtote-backend/library/database"
)

type CharacterRepository struct {
	*sqlx.DB
}

func (db *CharacterRepository) ExecTx(f func(tx *sqlx.Tx) error) error {
	return database.ExecTx(db, f)
}

func NewCharacterRepository(db *sqlx.DB) *CharacterRepository {
	return &CharacterRepository{db}
}
