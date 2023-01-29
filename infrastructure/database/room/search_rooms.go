package room

import (
	"encoding/json"
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
		order = "room_id"
	case "latest-post":
		order = "updated_at"
	case "posts-per-day":
		order = "posts_per_day"
	default:
		return nil, false, errors.New("順序指定が不正です")
	}

	inputs := make(map[string]any, 4)

	query := `
		WITH user_follows AS (
			SELECT followed FROM follows WHERE follower = :user
		)

		SELECT
			target_rooms.room_id,
			target_rooms.master_id,
			target_rooms.nickname,
			target_rooms.title,
			target_rooms.summary,
			target_rooms.tags,
			target_rooms.official,
			target_rooms.messages_count,
			target_rooms.members_count,
			target_rooms.updated_at,
			target_rooms.posts_per_day,
			COALESCE(
				JSON_AGG(JSON_BUILD_OBJECT(
					'id',       characters.id,
					'name',     characters.name,
					'mainicon', characters.mainicon
				)) FILTER (WHERE characters.name IS NOT NULL),
				'[]'
			)
		FROM
			(
				SELECT
					rooms.id AS room_id,
					characters.id AS master_id,
					characters.nickname,
					rooms.title,
					rooms.summary,
					ARRAY_REMOVE(ARRAY_AGG(rooms_tags.tag ORDER BY rooms_tags.id), NULL) AS tags,
					rooms.official,
					rooms.messages_count,
					rooms.members_count,
					rooms.updated_at,
					rooms.messages_count * 86400000.0 / GREATEST((EXTRACT(epoch from CURRENT_TIMESTAMP) - EXTRACT(epoch from rooms.created_at)), 259200000) AS posts_per_day
				FROM
					rooms
				LEFT JOIN
					rooms_tags ON (rooms.id = rooms_tags.room)
				JOIN
					characters ON (rooms.master = characters.id)
				WHERE
	`
	inputs["user"] = characterId

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

	if options.Participant != nil && *options.Participant == "own" {
		query += `
			EXISTS (
				SELECT
					*
				FROM
					rooms_members
				WHERE
					rooms_members.room = rooms.id AND
					rooms_members.member = :user
			) AND
		`
	}

	if options.Participant != nil && *options.Participant == "follow" {
		query += `
			EXISTS (
				SELECT
					*
				FROM
					rooms_members
				WHERE
					rooms_members.room = rooms.id AND
					rooms_members.member IN (SELECT * FROM user_follows)
			) AND
		`
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
					rooms.updated_at,
					posts_per_day
				ORDER BY
					` + order + ` ` + options.Sort + `
				LIMIT
					` + strconv.Itoa(options.Limit+1) + `
				OFFSET
					` + strconv.Itoa(options.Offset) + `
			) target_rooms
		LEFT JOIN
			rooms_members ON (target_rooms.room_id = rooms_members.room)
		LEFT JOIN
			characters ON (rooms_members.member = characters.id AND characters.id IN (SELECT * FROM user_follows) AND characters.deleted_at IS NULL)
		GROUP BY
			target_rooms.room_id,
			target_rooms.master_id,
			target_rooms.nickname,
			target_rooms.title,
			target_rooms.summary,
			target_rooms.tags,
			target_rooms.official,
			target_rooms.messages_count,
			target_rooms.members_count,
			target_rooms.updated_at,
			target_rooms.posts_per_day
		ORDER BY
			` + order + ` ` + options.Sort + `;
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
		var followedMembersJson string
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
			&room.PostsPerDay,
			&followedMembersJson,
		)

		err = json.Unmarshal([]byte(followedMembersJson), &room.FollowedMembers)
		if err != nil {
			return nil, false, err
		}

		rooms = append(rooms, room)
	}

	return
}
