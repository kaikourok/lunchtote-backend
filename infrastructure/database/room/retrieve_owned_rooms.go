package room

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/lib/pq"
)

func (db *RoomRepository) RetrieveOwnedRooms(characterId int) (rooms *[]model.RoomListItem, err error) {
	rows, err := db.Queryx(`
		SELECT
			rooms.id,
			characters.id,
			characters.nickname,
			rooms.title,
			rooms.summary,
			ARRAY_REMOVE(ARRAY_AGG(rooms_tags.tag ORDER BY rooms_tags.id), NULL),
			rooms.official,
			rooms.messages_count,
			rooms.members_count,
			rooms.updated_at
		FROM
			rooms
		JOIN
			characters ON (rooms.master = characters.id)
		LEFT JOIN
			rooms_tags ON (rooms.id = rooms_tags.room)
		WHERE
			rooms.master = $1 AND
			rooms.deleted_at IS NULL
		GROUP BY
			rooms.id,
			characters.id,
			characters.nickname,
			rooms.title,
			rooms.summary,
			rooms.official,
			rooms.messages_count,
			rooms.members_count,
			rooms.updated_at
		ORDER BY
			rooms.id;
	`, characterId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roomsSlice := make([]model.RoomListItem, 0, 64)

	for rows.Next() {
		var room model.RoomListItem
		err = rows.Scan(
			&room.Id,
			&room.Master.Id,
			&room.Master.Name,
			&room.Title,
			&room.Summary,
			pq.Array(&room.Tags),
			&room.Official,
			&room.MessagesCount,
			&room.MembersCount,
			&room.LastUpdate,
		)

		if err != nil {
			return nil, err
		}

		roomsSlice = append(roomsSlice, room)
	}

	return &roomsSlice, nil
}
