package repository

import "github.com/kaikourok/lunchtote-backend/entity/model"

type roomRepository interface {
	// ルーム取得
	// TODO: SearchRoom
	RetrieveOwnedRooms(characterId int) (rooms *[]model.RoomListItem, err error)

	// ルーム内処理関連
	RetrieveRoomOwnPermissions(characterId int, roomId int) (permissions *model.RoomMemberPermission, banned bool, err error)
	RetrieveRoomRelations(roomId int) (relations *model.RoomRelations, err error)
	RetrieveRoomMessages(characterId int, options *model.RoomMessageRetrieveSettings) (messages *[]model.RoomMessage, isContinueFollowing, isContinuePrevious *bool, err error)
	RetrieveRoomDetailData(characterId int, roomId int) (room *model.RoomDetailData, err error)
	PostRoomMessage(characterId int, message *model.RoomPostMessage, uploadPath string) error

	// ルーム関連情報取得
	RetrieveRoomTitle(roomId int) (title string, err error)

	// ルーム作成・削除関連
	CreateRoom(characterId int, room *model.Room) (roomId int, err error)
	DeleteRoom(characterId, roomId int) error
	RetrieveChildrenCreatableRooms(characterId int) (rooms *[]model.RoomOverview, err error)

	// 設定関連
	RetrieveRoomRoleSettings(roomId int) (roles *[]model.RoomRoleWithMembers, master int, title string, err error)
	RetrieveRoomGeneralSettings(roomId int) (room *model.Room, masterCharacter int, err error)

	// メンバー関連
	RetrieveRoomMembers(userId, roomId int) (members *[]model.RomeMemberWithRoles, err error)

	// BAN関連
	BanCharacterFromRoom(userId, targetId, roomId int) error
	CancelBanCharacterFromRoom(userId, targetId, roomId int) error
	RetrieveRoomBanStates(roomId int) (states *[]model.RoomBanState, master int, err error)

	// 招待関連
	InviteCharacterToRoom(userId, targetId, roomId int) error
	CancelInviteCharacterToRoom(userId, targetId, roomId int) error
	RetrieveRoomInviteStates(roomId int) (states *[]model.RoomInviteState, master int, err error)

	// ロール設定関連
	CreateRole(characterId, roomId int, roleName string, role *model.RoomRolePermission) (roleId int, err error)
	DeleteRole(characterId, roleId int) error
	UpdateRolePermissions(characterId, roleId int, roleName string, role *model.RoomRolePermission) error
	UpdateRolePriorities(characterId, roomId int, priorities *[]model.RoomRolePriority) error
	RevokeRoomRole(userId, targetId, role int) error
	GrantRoomRole(userId, targetId, roleId int) error
}
