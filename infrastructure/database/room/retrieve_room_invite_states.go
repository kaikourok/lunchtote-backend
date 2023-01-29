package room

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *RoomRepository) RetrieveRoomInviteStates(roomId int) (states *[]model.RoomInviteState, err error) {
	rows, err := db.Queryx(`
		SELECT
			invited_character.id,
			invited_character.name,
			invited_character.mainicon,
			inviter_character.id,
			inviter_character.name,
			inviter_character.mainicon,
			rooms_invited_characters.invited_at
		FROM
			rooms_invited_characters
		JOIN
			characters AS invited_character ON (rooms_invited_characters.invited = invited_character.id)
		JOIN
			characters AS inviter_character ON (rooms_invited_characters.inviter = inviter_character.id)
		WHERE
			rooms_invited_characters.room = $1
		ORDER BY
			rooms_invited_characters.invited_at DESC;
	`, roomId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statesSlice := make([]model.RoomInviteState, 0, 64)
	for rows.Next() {
		var state model.RoomInviteState
		err = rows.Scan(
			&state.Invited.Id,
			&state.Invited.Name,
			&state.Invited.Mainicon,
			&state.Inviter.Id,
			&state.Inviter.Name,
			&state.Inviter.Mainicon,
			&state.InvitedAt,
		)
		if err != nil {
			return nil, err
		}

		statesSlice = append(statesSlice, state)
	}

	return &statesSlice, nil
}
