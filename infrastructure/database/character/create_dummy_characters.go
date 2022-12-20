package character

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/library/secure"
)

func (db *CharacterRepository) CreateDummyCharacters(number int, name, nickname, summary, profile, password string, notificationTokenGenerator func() string) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		type InsertDummyCharacterStruct struct {
			Name              string `db:"name"`
			Nickname          string `db:"nickname"`
			Username          string `db:"username"`
			Summary           string `db:"summary"`
			Profile           string `db:"profile"`
			Password          string `db:"password"`
			NotificationToken string `db:"notification_token"`
		}

		salt, _ := uuid.NewUUID()

		characters := make([]InsertDummyCharacterStruct, number)
		for i := 0; i < len(characters); i++ {
			characters[i] = InsertDummyCharacterStruct{
				Name:              name,
				Nickname:          nickname,
				Username:          secure.GenerateShortHash(strconv.Itoa(i)+salt.String(), "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"),
				Summary:           summary,
				Profile:           profile,
				Password:          password,
				NotificationToken: notificationTokenGenerator(),
			}
		}

		_, err := tx.NamedExec(`
			INSERT INTO characters (
				name,
				nickname,
				username,
				summary,
				profile,
				password,
				notification_token
			) VALUES (
				:name,
				:nickname,
				:username,
				:summary,
				:profile,
				:password,
				:notification_token
			)
		`, characters)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE
				characters AS c
			SET
				id = (SELECT count(*) FROM characters WHERE characters.administrator = false AND characters.seq <= c.seq)
			WHERE
				c.administrator = false;
		`)

		return err
	})
}
