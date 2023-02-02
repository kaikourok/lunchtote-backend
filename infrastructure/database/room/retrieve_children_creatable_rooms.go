package room

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *RoomRepository) RetrieveChildrenCreatableRooms(characterId int) (rooms *[]model.RoomOverview, err error) {
	rows, err := db.Queryx(`
		SELECT
			rooms.id,
			rooms.title
		FROM
			rooms
		JOIN
			rooms_members ON (rooms.id = rooms_members.room)
		WHERE
			rooms_members.member = $1 AND
			rooms_members.create_children_room = true AND
			rooms.deleted_at IS NULL;
	`, characterId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roomsSlice := make([]model.RoomOverview, 0, 32)
	for rows.Next() {
		var room model.RoomOverview
		err := rows.Scan(&room.Id, &room.Title)
		if err != nil {
			return nil, err
		}

		roomsSlice = append(roomsSlice, room)
	}

	return &roomsSlice, nil
}
