package room

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *RoomRepository) RetrieveRoomMessageFetchConfig(characterId int) (configs *[]model.RoomMessageFetchConfigWithId, err error) {
	rows, err := db.Queryx(`
		SELECT
			id,
			name,
			room,
			search,
			refer_root,
			list,
			character,
			relate_filter,
			children,
			category
		FROM
			message_fetch_configs
		WHERE
			master = $1
		ORDER BY
			config_order;
	`, characterId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	configsSlice := make([]model.RoomMessageFetchConfigWithId, 0, 32)
	for rows.Next() {
		var config model.RoomMessageFetchConfigWithId
		err = rows.Scan(
			&config.Id,
			&config.Name,
			&config.Room,
			&config.Search,
			&config.ReferRoot,
			&config.List,
			&config.Character,
			&config.RelateFilter,
			&config.Children,
			&config.Category,
		)

		configsSlice = append(configsSlice, config)
	}

	return &configsSlice, nil
}
