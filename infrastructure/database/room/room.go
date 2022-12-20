package room

import (
	"github.com/jmoiron/sqlx"

	"github.com/kaikourok/lunchtote-backend/library/database"
)

type RoomRepository struct {
	*sqlx.DB
}

func (r *RoomRepository) ExecTx(f func(tx *sqlx.Tx) error) error {
	return database.ExecTx(r, f)
}

func NewRoomRepository(db *sqlx.DB) *RoomRepository {
	return &RoomRepository{db}
}
