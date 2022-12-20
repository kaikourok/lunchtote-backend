package repository

import (
	"bytes"

	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/entity/service"
)

type characterRepository interface {
	// キャラクター作成関連
	CreateCharacter(name, nickname, username, password, notificationToken string) (id int, err error)
	CreateAdministrator(id int, hashedPassword, name, nickname, username, notificationToken string) error
	CreateDummyCharacters(number int, name, nickname, summary, profile, password string, notificationTokenGenerator func() string) error

	// キャラクター情報取得関連
	RetrieveCharacterList(id *int, start, end int) (list *[]model.CharacterListItem, maxId int, err error)
	RetrieveProfile(userId *int, targetId int) (*model.Profile, error)
	RetrieveHomeData(userId int) (*model.HomeData, error)
	RetrieveCharacterNickname(characterId int) (nickname string, err error)

	// ログイン・登録情報関連
	RetrieveInitialData(id int) (existsUnreadNotification, existsUnreadMail bool, err error)
	RetrieveCredentials(id int) (password, notificationToken string, isAdministrator bool, err error)
	RetrievePassword(id int) (password string, err error)
	ExchangeEmailToId(email string) (id int, err error)
	ExchangeUsernameToId(username string) (id int, err error)
	ExchangeNotificationTokenToId(token string) (id int, err error)
	UpdatePasswordByResetCode(id int, code, password string) error
	UpdatePassword(id int, password string) error
	SetConfirmCode(id int, email, code string, expireMinutes int) error
	SetPasswordResetCode(id int, email, code string, expireMinutes int) error
	RetrieveConfirmCode(id int) (code, email string, err error)
	DeleteCharacter(id int) error
	UndeleteCharacter(id int) error
	CheckUsernameExists(username string) (exists bool, err error)

	//SSO関連
	RetrieveCredentialsByTwitter(twitterId string) (characterId int, notificationToken string, err error)
	RetrieveCredentialsByGoogle(googleId string) (characterId int, notificationToken string, err error)
	RegisterGoogleData(characterId int, googleId string) error
	RegisterTwitterData(characterId int, twitterId string) error

	// プロフィール編集関連
	RetrieveProfileEditData(id int) (data *model.ProfileEditData, err error)
	UpdateProfile(id int, profile *model.ProfileEditData) error
	RetrieveCharacterIcons(id int) (icons *[]model.Icon, err error)
	RetrieveCharacterProfileImages(id int) (images *[]model.ProfileImage, err error)
	UpdateIcons(id int, icons *[]model.Icon, insertOnly bool) error
	UpdateProfileImages(id int, images *[]model.ProfileImage) error
	UpdateEmail(id int, email string) error

	// 関連性関連
	Follow(userId, targetId int) (userName string, webhook string, err error)
	Mute(userId, targetId int) error
	Block(userId, targetId int) error
	Unmute(userId, targetId int) error
	Unfollow(userId, targetId int) error
	Unblock(userId, targetId int) error
	RetrieveFollowList(userId, targetId int) (list *[]model.CharacterListItem, err error)
	RetrieveFollowerList(userId, targetId int) (list *[]model.CharacterListItem, err error)
	RetrieveRelatedFollowerList(userId, targetId int) (list *[]model.CharacterListItem, err error)
	RetrieveMuteList(id int) (list *[]model.CharacterListItem, err error)
	RetrieveBlockList(id int) (list *[]model.CharacterListItem, err error)

	// リスト操作
	CreateList(characterId int, name string) (listId int, err error)
	DeleteList(userId, listId int) error
	AddCharacterToList(userId, targetId, listId int) error
	RemoveCharacterFromList(userId, targetId, listId int) error
	RetrieveLists(id int) (lists *[]model.ListOverview, err error)

	// 画像管理関連
	RetrieveUploadedImages(id int) (images *[]model.UploadedImage, err error)
	DeleteUploadedImages(characterId int, imageIds *[]int, uploadDir string) error
	SaveUploadImageData(id int, imageBuffers []*bytes.Buffer, imageType service.ImageTypeId, uploadDir string) (*[]string, error)
	CreateLayeringGroup(characterId int, name string) (id int, err error)
	DeleteLayeringGroup(characterId, groupId int) error
	UpdateLayeringGroupName(characterId, groupId int, name string) error
	CreateLayerGroup(characterId, groupId int, name string) (id int, err error)
	DeleteLayerGroup(characterId, groupId int) error
	UpdateLayerGroupName(characterId, groupId int, name string) error
	UpdateLayerItems(characterId, groupId int, items *[]model.CharacterIconLayerItemEditData) (result *[]model.CharacterIconLayerItem, err error)
	RetrieveLayeringGroupOverviews(characterId int) (overviews *[]model.CharacterIconLayeringGroupOverview, err error)
	RetrieveLayeringGroup(characterId, groupId int) (layeringGroup *model.CharacterIconLayeringGroup, err error)
	DeleteLayerItems(characterId int, itemIds *[]int) error
	CreateProcessSchema(characterId, groupId int, process *model.CharacterIconProcessSchema) (id int, err error)
	DeleteProcessSchema(characterId, processId int) error
	UpdateLayerGroupOrders(characterId, groupId int, orders *[]model.CharacterIconLayerGroupOrder) error

	// サジェスト関連
	RetrieveCharacterSuggestions(id int, searchText string, excludeOwn bool) (suggestions *model.CharacterSuggestionsData, err error)
	RetrieveListSuggestions(characterId int, searchText string, listId int) (suggestions *model.CharacterSuggestionsData, err error)

	// BAN関連
	Adomonish(targetId int, data *model.AdomonishData) error
	Ban(targetId int, data *model.BanData) error
	Unban(targetId int, data *model.UnbanData) error
	RetrieveProhibitionRelatedData(targetId int) (data *[]model.ProhibitionRelatedData, err error)

	// その他
	RetrieveNotifications(id, start, number int) (notifications *[]model.Notification, isContinue bool, err error)
	RetrieveEmailRegistratedCharacters(targetCharacters *[]int) (characters *[]model.CharacterEmailRegistratedData, err error)
}
