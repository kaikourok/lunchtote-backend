package room

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *RoomRepository) AddRoomMessageFetchConfig(characterId int, config *model.RoomMessageFetchConfig) (configId int, err error) {
	row := db.QueryRowx(`
		INSERT INTO message_fetch_configs (
			master,
			config_order,
			name,
			room,
			search,
			refer_root,
			list,
			character,
			relate_filter,
			children,
			category
		) VALUES (
			$1,
			COALESCE((SELECT config_order + 1 FROM message_fetch_configs WHERE master = $1 ORDER BY config_order DESC LIMIT 1), 0),
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10
		)

		RETURNING
			id;
	`,
		characterId,
		config.Name,
		config.Room,
		config.Search,
		config.ReferRoot,
		config.List,
		config.Character,
		config.RelateFilter,
		config.Children,
		config.Category,
	)

	err = row.Scan(&configId)
	return
}
