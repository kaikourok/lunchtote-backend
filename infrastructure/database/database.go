package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/infrastructure/database/character"
	"github.com/kaikourok/lunchtote-backend/infrastructure/database/diary"
	"github.com/kaikourok/lunchtote-backend/infrastructure/database/forum"
	"github.com/kaikourok/lunchtote-backend/infrastructure/database/general"
	"github.com/kaikourok/lunchtote-backend/infrastructure/database/mail"
	"github.com/kaikourok/lunchtote-backend/infrastructure/database/room"
)

type PostgresRepository struct {
	character.CharacterRepository
	diary.DiaryRepository
	general.GeneralRepository
	mail.MailRepository
	room.RoomRepository
	forum.ForumRepository
}

func NewRepository(dsn *DataSource) (*PostgresRepository, error) {
	db, err := sqlx.Connect("postgres", dsn.getDataSourceString(
		dsnOptionSslMode("disable"),
	))
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{
		*character.NewCharacterRepository(db),
		*diary.NewDiaryRepository(db),
		*general.NewGeneralRepository(db),
		*mail.NewMailRepository(db),
		*room.NewRoomRepository(db),
		*forum.NewForumRepository(db),
	}, nil
}
