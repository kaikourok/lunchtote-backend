package room

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *RoomRepository) RetrieveRoomBanStates(roomId int) (states *[]model.RoomBanState, master int, err error) {
	rows, err := db.Queryx(`
		SELECT
			(SELECT master FROM rooms WHERE id = $1),
			banned_character.id,
			banned_character.name,
			banned_character.mainicon,
			banner_character.id,
			banner_character.name,
			banner_character.mainicon,
			rooms_banned_character.banned_at
		FROM
			rooms_banned_characters
		JOIN
			characters AS banned_character ON (rooms_banned_characters.banned = banned_character.id)
		JOIN
			characters AS banner_character ON (rooms_banned_characters.banner = banner_character.id)
		WHERE
			rooms_banned_characters.room = $1
		ORDER BY
			rooms_banned_characters.banned_at DESC;
	`)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	statesSlice := make([]model.RoomBanState, 0, 64)

	for rows.Next() {
		var state model.RoomBanState
		err = rows.Scan(
			&state.Banned.Id,
			&state.Banned.Name,
			&state.Banned.Mainicon,
			&state.Banner.Id,
			&state.Banner.Name,
			&state.Banner.Mainicon,
			&state.BannedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		statesSlice = append(statesSlice, state)
	}

	return &statesSlice, 0, nil
}
