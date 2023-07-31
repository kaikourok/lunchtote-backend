package repository

import "github.com/kaikourok/lunchtote-backend/entity/model"

type roomRepository interface {
	// ルーム取得
	RetrieveOwnedRooms(characterId int) (rooms *[]model.RoomListItem, err error)
	RetrieveMemberRooms(characterId int) (rooms *[]model.RoomListItem, err error)
	RetrieveInvitedRooms(characterId int) (rooms *[]model.RoomListItem, err error)
	SearchRooms(characterId int, options *model.RoomSearchOptions) (rooms []model.RoomListItem, isContinue bool, err error)

	// ルーム内処理関連
	RetrieveRoomOwnPermissions(characterId int, roomId int) (permissions *model.RoomMemberPermission, roleType string, banned bool, err error)
	RetrieveRoomRelations(roomId int) (relations *model.RoomRelations, err error)
	RetrieveRoomMessages(characterId int, options *model.RoomMessageRetrieveSettings) (messages *[]model.RoomMessage, isContinueFollowing, isContinuePrevious *bool, err error)
	RetrieveRoomDetailData(characterId int, roomId int) (room *model.RoomDetailData, err error)
	PostRoomMessage(characterId int, message *model.RoomPostMessage, uploadPath string) (messageId int, err error)

	// 通知関連
	NotificateRoomMessage(messageId int) (dto *model.RoomNotificationRelatedData, err error)
	RetrieveRoomSubscribeStates(characterId, roomId int) (states *model.RoomSubscribeStates, err error)
	SubscribeRoomMessage(characterId, roomId int) error
	SubscribeRoomNewMember(characterId, roomId int) error
	UnsubscribeRoomMessage(characterId, roomId int) error
	UnsubscribeRoomNewMember(characterId, roomId int) error

	// メッセージ取得関連
	AddRoomMessageFetchConfig(characterId int, config *model.RoomMessageFetchConfig) (configId int, err error)
	DeleteRoomMessageFetchConfig(characterId, configId int) error
	RenameRoomMessageFetchConfig(characterId, configId int, name string) error
	RetrieveRoomMessageFetchConfig(characterId int) (configs *[]model.RoomMessageFetchConfigWithId, err error)
	UpdateRoomMessageFetchConfigOrders(characterId int, orders *[]model.RoomMessageFetchConfigOrder) error

	// ルーム関連情報取得
	RetrieveRoomTitle(roomId int) (title string, err error)

	// ルーム作成・削除関連
	CreateRoom(characterId int, room *model.Room) (roomId int, err error)
	DeleteRoom(characterId, roomId int) error
	RetrieveChildrenCreatableRooms(characterId int) (rooms *[]model.RoomOverview, err error)

	// 設定関連
	RetrieveRoomRoleSettings(roomId int) (roles []model.RoomRole, master int, err error)
	RetrieveRoomGeneralSettings(roomId int) (room *model.Room, masterCharacter int, err error)
	UpdateRoomSettings(characterId, roomId int, room *model.Room) error

	// メンバー関連
	RetrieveRoomMembers(userId, roomId int) (members *[]model.RomeMemberWithRoles, err error)

	// BAN関連
	BanCharacterFromRoom(userId, targetId, roomId int) error
	CancelBanCharacterFromRoom(userId, targetId, roomId int) error
	RetrieveRoomBanStates(roomId int) (states *[]model.RoomBanState, err error)

	// 招待関連
	InviteCharacterToRoom(userId, targetId, roomId int) error
	CancelInviteCharacterToRoom(userId, targetId, roomId int) error
	RetrieveRoomInviteStates(roomId int) (states *[]model.RoomInviteState, err error)
	RetrieveRoomInviteSuggestions(characterId int, searchText string, roomId int) (suggestions *model.CharacterSuggestionsData, err error)

	// ロール設定関連
	JoinToRoom(targetId, roomId int) (room *model.RoomOverview, targetName string, newMemberWebhooks []string, err error)
	CreateRole(characterId, roomId int, roleName string, role *model.RoomRolePermission) (roleId int, err error)
	DeleteRole(characterId, roleId int) error
	UpdateRolePermissions(characterId, roleId int, roleName string, role *model.RoomRolePermission) error
	UpdateRolePriorities(characterId, roomId int, priorities *[]model.RoomRolePriority) error
	RevokeRoomRole(userId, targetId, role int) error
	GrantRoomRole(userId, targetId, roleId int) error
}
