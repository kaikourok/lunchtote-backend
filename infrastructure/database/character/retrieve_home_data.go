package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) RetrieveHomeData(userId int) (home *model.HomeData, err error) {
	row := db.QueryRowx(`
		SELECT
			(
				SELECT
					nickname
				FROM
					characters
				WHERE
					id = $1
			);
	`, userId)
	home = &model.HomeData{}
	err = row.Scan(
		&home.Nickname,
	)
	if err != nil {
		return nil, err
	}

	return home, nil
}
