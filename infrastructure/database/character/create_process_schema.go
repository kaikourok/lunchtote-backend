package character

import (
	"errors"

	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *CharacterRepository) CreateProcessSchema(characterId, groupId int, process *model.CharacterIconProcessSchema) (id int, err error) {
	var master int
	row := db.QueryRowx(`
		SELECT
			character
		FROM
			characters_icon_layering_groups
		WHERE
			id = $1;
	`, groupId)

	err = row.Scan(&master)
	if err != nil {
		return 0, err
	}

	if characterId != master {
		return 0, errors.New("レイヤリンググループの所有者ではありません")
	}

	row = db.QueryRowx(`
		INSERT INTO characters_icon_process_schemas (
			layering_group,
			name,
			x,
			y,
			scale,
			rotate
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
		RETURNING
			id;
	`,
		groupId,
		process.Name,
		process.X,
		process.Y,
		process.Scale,
		process.Rotate,
	)

	err = row.Scan(&id)
	return
}
