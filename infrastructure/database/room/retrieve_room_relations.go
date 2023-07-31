package room

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *RoomRepository) RetrieveRoomRelations(roomId int) (room *model.RoomRelations, err error) {
	rows, err := db.Queryx(`
		WITH room AS (
			SELECT
				title,
				belong
			FROM
				rooms
			WHERE
				id = $1
		), parent_room_children_referable AS (
			SELECT 
				EXISTS (
					SELECT
						*
					FROM
						rooms
					WHERE
						id = (SELECT belong FROM room) AND
						children_referable = true
				)
		), results AS (
			SELECT
				'parent' AS type,
				(
					SELECT
						belong
					FROM
						room
				) AS id,
				(
					SELECT
						title
					FROM
						rooms
					WHERE
						id = (SELECT belong FROM room)
				) AS title

			UNION ALL

			SELECT
				'sibling' AS type,
				$1 AS id,
				title
			FROM
				room

			UNION ALL

			SELECT
				'sibling' AS type,
				id,
				title
			FROM
				rooms
			WHERE
				(SELECT * FROM parent_room_children_referable) = true AND
				id != $1 AND
				deleted_at IS NULL AND
				belong = (SELECT belong FROM room)

			UNION ALL

			SELECT
				'children' AS type,
				id,
				title
			FROM
				rooms
			WHERE
				belong = $1 AND
				deleted_at IS NULL
		)
			
		SELECT
			*
		FROM
			results
		ORDER BY
			id;
	`, roomId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parent *model.RoomOverview
	siblings := make([]model.RoomOverview, 0, 16)
	children := make([]model.RoomOverview, 0, 16)

	for rows.Next() {
		var rowType string
		var roomId *int
		var roomTitle *string
		err := rows.Scan(
			&rowType,
			&roomId,
			&roomTitle,
		)
		if err != nil {
			return nil, err
		}

		if roomId == nil || roomTitle == nil {
			continue
		}

		room := model.RoomOverview{
			Id:    *roomId,
			Title: *roomTitle,
		}

		switch rowType {
		case "parent":
			parent = &room
		case "sibling":
			siblings = append(siblings, room)
		case "children":
			children = append(children, room)
		}
	}

	return &model.RoomRelations{
		Parent:   parent,
		Siblings: siblings,
		Children: children,
	}, nil
}
