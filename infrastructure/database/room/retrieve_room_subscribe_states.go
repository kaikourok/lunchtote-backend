package room

import (
	model "github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *RoomRepository) RetrieveRoomSubscribeStates(characterId, roomId int) (states *model.RoomSubscribeStates, err error) {
	row := db.QueryRowx(`
		SELECT
			EXISTS (SELECT * FROM rooms_message_subscribers    WHERE character = $1 AND room = $2),
			EXISTS (SELECT * FROM rooms_new_member_subscribers WHERE character = $1 AND room = $2);
	`, characterId, roomId)

	states = &model.RoomSubscribeStates{}
	err = row.Scan(&states.Message, &states.NewMember)

	return
}
