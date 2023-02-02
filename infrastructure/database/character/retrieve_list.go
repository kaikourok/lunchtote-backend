package character

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"golang.org/x/sync/errgroup"
)

func (db *CharacterRepository) RetrieveList(listId int) (listName string, characters []model.CharacterOverview, err error) {
	var eg errgroup.Group

	eg.Go(func() error {
		row := db.QueryRowx(`
			SELECT
				name
			FROM
				lists
			WHERE
				id = $1;
		`, listId)

		return row.Scan(&listName)
	})

	eg.Go(func() error {
		rows, err := db.Queryx(`
			SELECT
				characters.id,
				characters.nickname,
				characters.mainicon
			FROM
				lists_characters
			JOIN
				characters ON (characters.id = lists_characters.character)
			WHERE
				list = $1
			ORDER BY
				characters.id;
		`, listId)
		if err != nil {
			return err
		}
		defer rows.Close()

		characters = make([]model.CharacterOverview, 0, 64)
		for rows.Next() {
			var character model.CharacterOverview
			err = rows.Scan(
				&character.Id,
				&character.Name,
				&character.Mainicon,
			)
			if err != nil {
				return err
			}

			characters = append(characters, character)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return "", nil, err
	}

	return
}
