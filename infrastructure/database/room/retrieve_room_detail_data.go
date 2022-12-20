package room

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/lib/pq"
)

func (db *RoomRepository) RetrieveRoomDetailData(characterId int, roomId int) (room *model.RoomDetailData, err error) {
	room = &model.RoomDetailData{}

	row := db.QueryRowx(`
		SELECT
			rooms.title,
			rooms.summary,
			rooms.description,
			rooms.searchable,
			rooms.allow_recommendation,
			rooms.official,
			rooms.members_count,
			rooms.updated_at,
			characters.id,
			characters.name,
			characters.mainicon,
			ARRAY_REMOVE(ARRAY_AGG(rooms_tags.tag ORDER BY rooms_tags.id), NULL)
		FROM 
			rooms
		JOIN
			characters ON (rooms.master = characters.id)
		LEFT JOIN
			rooms_tags ON (rooms.id = rooms_tags.room)
		WHERE
			rooms.id = $1
		GROUP BY
			rooms.title,
			rooms.summary,
			rooms.description,
			rooms.searchable,
			rooms.allow_recommendation,
			rooms.official,
			rooms.members_count,
			rooms.updated_at,
			characters.id,
			characters.name,
			characters.mainicon;
	`, roomId)

	err = row.Scan(
		&room.Title,
		&room.Summary,
		&room.Description,
		&room.Searchable,
		&room.AllowRecommendation,
		&room.Official,
		&room.MembersCount,
		&room.UpdatedAt,
		&room.Master.Id,
		&room.Master.Name,
		&room.Master.Mainicon,
		pq.Array(&room.Tags),
	)
	return
}
