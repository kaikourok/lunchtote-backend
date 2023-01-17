package room

import (
	"errors"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/lib/pq"
)

func (db *RoomRepository) SearchRooms(characterId int, options *model.RoomSearchOptions) (rooms []model.RoomListItem, isContinue bool, err error) {
	var order string
	switch options.Order {
	case "id":
		order = "id"
	case "latest-post":
		order = "updated_at"
	default:
		return nil, false, errors.New("順序指定が不正です")
	}

	inputs := make(map[string]any, 4)

	query := `
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
		LEFT JOIN
			rooms_tags ON (rooms.id = rooms_tags.room)
		JOIN
			characters ON (rooms.master = characters.id)
		WHERE
	`

	if options.Title != "" {
		query += `
			rooms.title LIKE likequery(:title) AND
		`
		inputs["title"] = options.Title
	}

	if options.Description != "" {
		query += `
			rooms.description LIKE likequery(:description) AND
		`
		inputs["description"] = options.Description
	}

	if 0 < len(options.Tags) {
		query += `
			EXISTS (
				SELECT 
					*
				FROM
					rooms_tags
				WHERE
					rooms_tags.room = rooms.id AND
					rooms_tags.tag IN :tags
			) AND
		`
		inputs["tags"] = options.Tags
	}

	if 0 < len(options.ExcludedTags) {
		query += `
			NOT EXISTS (
				SELECT 
					*
				FROM
					rooms_tags
				WHERE
					rooms_tags.room = rooms.id AND
					rooms_tags.tag IN :excluded_tags
			) AND
		`
		inputs["excluded_tags"] = options.ExcludedTags
	}

	query += `
			rooms.deleted_at IS NULL AND
			rooms.searchable = true
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
			` + order + ` ` + options.Sort + `
		LIMIT
			` + strconv.Itoa(options.Limit+1) + `
		OFFSET
			` + strconv.Itoa(options.Offset) + `;
	`

	query, args, err := sqlx.Named(query, inputs)
	if err != nil {
		return nil, false, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, false, err
	}

	rows, err := db.Queryx(db.Rebind(query), args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	rooms = make([]model.RoomListItem, 0, options.Limit)
	for i := 0; rows.Next(); i++ {
		if i == options.Limit {
			isContinue = true
			break
		}
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
		rooms = append(rooms, room)
	}

	return
}
