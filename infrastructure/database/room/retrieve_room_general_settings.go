package room

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/lib/pq"
)

func (db *RoomRepository) RetrieveRoomGeneralSettings(roomId int) (room *model.Room, masterCharacter int, err error) {
	row := db.QueryRowx(`
		SELECT
			rooms.master,
			rooms.title,
			rooms.summary,
			rooms.description,
			ARRAY_REMOVE(ARRAY_AGG(rooms_tags.tag ORDER BY rooms_tags.id), NULL),
			rooms.searchable,
			rooms.allow_recommendation
		FROM
			rooms
		LEFT JOIN
			rooms_tags ON (rooms.id = rooms_tags.room)
		WHERE
			rooms.id = $1 AND
			rooms.deleted_at IS NULL
		GROUP BY
			rooms.id,
			rooms.master,
			rooms.title,
			rooms.summary,
			rooms.description,
			rooms.searchable,
			rooms.allow_recommendation
		ORDER BY
			rooms.id;
	`, roomId)

	room = &model.Room{}
	err = row.Scan(
		&masterCharacter,
		&room.Title,
		&room.Summary,
		&room.Description,
		pq.Array(&room.Tags),
		&room.Searchable,
		&room.AllowRecommendation,
	)
	if err != nil {
		return nil, 0, err
	}

	return room, masterCharacter, nil
}
