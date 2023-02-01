package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	characterController "github.com/kaikourok/lunchtote-backend/cmd/server/controller/character"
	controlController "github.com/kaikourok/lunchtote-backend/cmd/server/controller/control"
	diaryController "github.com/kaikourok/lunchtote-backend/cmd/server/controller/diary"
	forumController "github.com/kaikourok/lunchtote-backend/cmd/server/controller/forum"
	generalController "github.com/kaikourok/lunchtote-backend/cmd/server/controller/general"
	mailController "github.com/kaikourok/lunchtote-backend/cmd/server/controller/mail"
	roomController "github.com/kaikourok/lunchtote-backend/cmd/server/controller/room"
	"github.com/kaikourok/lunchtote-backend/cmd/server/middleware"
	"github.com/kaikourok/lunchtote-backend/registry"
	"github.com/kaikourok/lunchtote-backend/usecase/character"
	"github.com/kaikourok/lunchtote-backend/usecase/control"
	"github.com/kaikourok/lunchtote-backend/usecase/diary"
	"github.com/kaikourok/lunchtote-backend/usecase/forum"
	"github.com/kaikourok/lunchtote-backend/usecase/general"
	"github.com/kaikourok/lunchtote-backend/usecase/mail"
	"github.com/kaikourok/lunchtote-backend/usecase/room"
)

func NewRouter(registry registry.Registry) *gin.Engine {
	config := registry.GetConfig()
	router := gin.New()

	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	store, err := redis.NewStore(
		10,
		"tcp",
		config.GetString("session.host")+":"+config.GetString("session.port"),
		config.GetString("session.password"),
		[]byte(config.GetString("session.secret")),
	)
	if err != nil {
		log.Fatal("Redisへの接続に失敗しました")
	}

	err = redis.SetKeyPrefix(store, config.GetString("session.redis-prefix"))
	if err != nil {
		log.Fatal("セッションストアの設定中にエラーが発生しました")
	}

	router.Use(sessions.Sessions(config.GetString("session.name"), store))

	general := generalController.NewGeneralController(general.NewGeneralUsecase(registry))
	character := characterController.NewCharacterController(character.NewCharacterUsecase(registry), registry)
	room := roomController.NewRoomController(room.NewRoomUsecase(registry))
	mail := mailController.NewMailController(mail.NewMailUsecase(registry))
	diary := diaryController.NewDiaryController(diary.NewDiaryUsecase(registry))
	forum := forumController.NewForumController(forum.NewForumUsecase(registry))
	control := controlController.NewControlController(control.NewControlUsecase(registry))

	{
		api := router.Group("")
		api.GET("/health", general.Status)
		api.GET("/connections", general.Connections)
		api.POST("/inquiry", general.Inquiry)
		api.POST("/signin", character.SignIn)
		api.POST("/signout", middleware.Auth(), character.SignOut)
		api.GET("/initial-data", middleware.Auth(), character.RetrieveInitialData)
		api.POST("/mail-confirm", middleware.Auth(), character.ConfirmEmail)
		api.POST("/reset-password", character.RequestPasswordResetCode)
		api.POST("/reset-password-confirm", character.UpdatePasswordByResetCode)
		api.GET("/announcements", general.RetrieveAnnouncementOverviews)
		api.GET("/announcements/:id", general.RetrieveAnnouncement)

		{
			oauthGroup := api.Group("oauth")
			oauthGroup.GET("/google", character.GoogleOauthRequest)
			oauthGroup.GET("/google/callback", character.GoogleOauthCallback)
			oauthGroup.POST("/google/unlink", middleware.Auth(), character.UnlinkGoogle)
			oauthGroup.GET("/twitter", character.TwitterOauthRequest)
			oauthGroup.GET("/twitter/callback", character.TwitterOauthCallback)
			oauthGroup.POST("/twitter/unlink", middleware.Auth(), character.UnlinkTwitter)
		}

		{
			charactersGroup := api.Group("characters")
			charactersGroup.POST("", character.SignUp)
			charactersGroup.GET("", character.RetrieveCharacterList)
			charactersGroup.POST("/exists", character.CheckUsernameExists)
			charactersGroup.POST("/inline-search", middleware.Auth(), character.RetrieveCharacterSuggestions)
			charactersGroup.POST("/exchange-notification-token", character.ExchangeNotificationTokenToId)

			{
				characterGroup := charactersGroup.Group(":id")
				characterGroup.GET("", character.RetrieveProfile)
				characterGroup.GET("/following", middleware.Auth(), character.RetrieveFollowList)
				characterGroup.GET("/followers", middleware.Auth(), character.RetrieveFollowerList)
				characterGroup.GET("/followers/related", middleware.Auth(), character.RetrieveRelatedFollowerList)
				characterGroup.POST("/follow", middleware.Auth(), character.Follow)
				characterGroup.POST("/mute", middleware.Auth(), character.Mute)
				characterGroup.POST("/block", middleware.Auth(), character.Block)
				characterGroup.POST("/unfollow", middleware.Auth(), character.Unfollow)
				characterGroup.POST("/unmute", middleware.Auth(), character.Unmute)
				characterGroup.POST("/unblock", middleware.Auth(), character.Unblock)
			}

			{
				mainGroup := charactersGroup.Group("main")
				mainGroup.GET("/home", middleware.Auth(), character.RetrieveHomeData)
				mainGroup.GET("/muting", middleware.Auth(), character.RetrieveMuteList)
				mainGroup.GET("/blocking", middleware.Auth(), character.RetrieveBlockList)
				mainGroup.POST("/upload", middleware.Auth(), character.UploadImages)
				mainGroup.POST("/upload/base64", middleware.Auth(), character.UploadBase64EncordedImages)
				mainGroup.GET("/notifications", middleware.Auth(), character.RetrieveNotifications)
				mainGroup.POST("/notifications/checked", middleware.Auth(), character.UpdateNotificationChecked)
				mainGroup.POST("/delete", middleware.Auth(), character.DeleteCharacter)

				{
					settingGroup := mainGroup.Group("settings")
					settingGroup.GET("/profile", middleware.Auth(), character.RetrieveProfileEditData)
					settingGroup.POST("/profile", middleware.Auth(), character.UpdateProfile)
					settingGroup.GET("/icons", middleware.Auth(), character.RetrieveCharacterIconsEditData)
					settingGroup.POST("/icons", middleware.Auth(), character.UpdateIcons)
					settingGroup.GET("/profile-images", middleware.Auth(), character.RetrieveCharacterProfileImagesEditData)
					settingGroup.POST("/profile-images", middleware.Auth(), character.UpdateProfileImages)
					settingGroup.GET("/uploaded-images", middleware.Auth(), character.RetrieveUploadedImages)
					settingGroup.POST("/uploaded-images/delete", middleware.Auth(), character.DeleteUploadedImages)
					settingGroup.GET("/other", middleware.Auth(), character.RetrieveOtherSettings)
					settingGroup.POST("/other", middleware.Auth(), character.UpdateOtherSettings)
					settingGroup.POST("/email", middleware.Auth(), character.RequestRegisterEmail)
					settingGroup.POST("/email/unregister", middleware.Auth(), character.UnregisterEmail)

					{
						layeringGroup := settingGroup.Group("layerings")
						layeringGroup.GET("", middleware.Auth(), character.RetrieveLayeringGroupOverviews)
						layeringGroup.POST("", middleware.Auth(), character.CreateLayeringGroup)
						layeringGroup.GET("/:layering_group", middleware.Auth(), character.RetrieveLayeringGroup)
						layeringGroup.PUT("/:layering_group", middleware.Auth(), character.UpdateLayeringGroupName)
						layeringGroup.POST("/:layering_group/delete", middleware.Auth(), character.DeleteLayeringGroup)
						layeringGroup.POST("/:layering_group/delete-items", middleware.Auth(), character.DeleteLayerItems)
						layeringGroup.POST("/:layering_group/layers", middleware.Auth(), character.CreateLayerGroup)
						layeringGroup.POST("/:layering_group/layer-orders", middleware.Auth(), character.UpdateLayerGroupOrders)
						layeringGroup.PUT("/:layering_group/layers/:layer_group", middleware.Auth(), character.UpdateLayerGroupName)
						layeringGroup.POST("/:layering_group/layers/:layer_group/delete", middleware.Auth(), character.DeleteLayerGroup)
						layeringGroup.PUT("/:layering_group/layers/:layer_group/images", middleware.Auth(), character.UpdateLayerItems)
						layeringGroup.POST("/:layering_group/processes", middleware.Auth(), character.CreateProcessSchema)
						layeringGroup.POST("/:layering_group/processes/:process/delete", middleware.Auth(), character.DeleteProcessSchema)
					}
				}

				{
					listGroup := mainGroup.Group("lists")
					listGroup.POST("", middleware.Auth(), character.CreateList)
					listGroup.GET("", middleware.Auth(), character.RetrieveLists)
					listGroup.GET("/:list", middleware.Auth(), character.RetrieveList)
					listGroup.POST("/:list/add", middleware.Auth(), character.AddCharacterToList)
					listGroup.POST("/:list/remove", middleware.Auth(), character.RemoveCharacterFromList)
					listGroup.POST("/:list/rename", middleware.Auth(), character.RenameList)
					listGroup.POST("/:list/delete", middleware.Auth(), character.DeleteList)
					listGroup.POST("/:list/search-target", middleware.Auth(), character.RetrieveListSuggestions)
				}
			}
		}

		{
			roomsGroup := api.Group("rooms")
			roomsGroup.POST("", middleware.Auth(), room.CreateRoom)
			roomsGroup.GET("/create-data", middleware.Auth(), room.RetrieveRoomCreateData)
			roomsGroup.GET("/owned", middleware.Auth(), room.RetrieveOwnedRooms)
			roomsGroup.GET("/membered", middleware.Auth(), room.RetrieveMemberRooms)
			roomsGroup.GET("/messages", middleware.Auth(), room.RetrieveRoomMessages)
			roomsGroup.GET("/message-edit-data", middleware.Auth(), room.RetrieveRoomMessageEditRequiredData)
			roomsGroup.POST("/search", middleware.Auth(), room.SearchRooms)

			{
				fetchConfigGroup := roomsGroup.Group("fetch-configs")
				fetchConfigGroup.GET("", middleware.Auth(), room.RetrieveRoomMessageFetchConfig)
				fetchConfigGroup.POST("/add", middleware.Auth(), room.AddRoomMessageFetchConfig)
				fetchConfigGroup.POST("/rename", middleware.Auth(), room.RenameRoomMessageFetchConfig)
				fetchConfigGroup.POST("/delete", middleware.Auth(), room.DeleteRoomMessageFetchConfig)
				fetchConfigGroup.POST("/orders", middleware.Auth(), room.UpdateRoomMessageFetchConfigOrders)
			}

			{
				roomGroup := roomsGroup.Group(":id")
				roomGroup.GET("", middleware.Auth(), room.RetrieveRoomInitialData)
				roomGroup.GET("/permissions", middleware.Auth(), room.RetrieveRoomOwnPermissions)
				roomGroup.POST("/messages", middleware.Auth(), room.PostRoomMessage)
				roomGroup.POST("/message-event/subscribe", middleware.Auth(), room.SubscribeRoomMessage)
				roomGroup.POST("/message-event/unsubscribe", middleware.Auth(), room.UnsubscribeRoomMessage)
				roomGroup.POST("/new-member-event/subscribe", middleware.Auth(), room.SubscribeRoomNewMember)
				roomGroup.POST("/new-member-event/unsubscribe", middleware.Auth(), room.UnsubscribeRoomNewMember)

				{
					controlGroup := roomGroup.Group("control")
					controlGroup.GET("/general", middleware.Auth(), room.RetrieveRoomGeneralSettings)
					controlGroup.GET("/members", middleware.Auth(), room.RetrieveRoomMembers)
					controlGroup.GET("/invite", middleware.Auth(), room.RetrieveRoomInviteStates)
					controlGroup.GET("/ban", middleware.Auth(), room.RetrieveRoomBanStates)
					controlGroup.POST("/grant-role", middleware.Auth(), room.GrantRoomRole)
					controlGroup.POST("/revoke-role", middleware.Auth(), room.RevokeRoomRole)
					controlGroup.POST("/invite", middleware.Auth(), room.InviteCharacterToRoom)
					controlGroup.POST("/ban", middleware.Auth(), room.BanCharacter)
					controlGroup.POST("/cancel-invite", middleware.Auth(), room.CancelInviteCharacterToRoom)
					controlGroup.POST("/cancel-ban", middleware.Auth(), room.CancelBanCharacterFromRoom)
					controlGroup.POST("/delete", middleware.Auth(), room.DeleteRoom)
					controlGroup.POST("/search-invite-target", middleware.Auth(), room.RetrieveRoomInviteSuggestions)

					{
						roleGroup := controlGroup.Group("role")
						roleGroup.GET("", middleware.Auth(), room.RetrieveRoomRoleSettings)
						roleGroup.POST("/create", middleware.Auth(), room.CreateRole)
						roleGroup.POST("/delete", middleware.Auth(), room.DeleteRole)
						roleGroup.POST("/:role/update", middleware.Auth(), room.UpdateRolePermissions)
						roleGroup.POST("/update-priorities", middleware.Auth(), room.UpdateRolePriorities)
					}
				}
			}
		}

		{
			mailGroup := api.Group("mails")
			mailGroup.GET("", middleware.Auth(), mail.RetrieveMails)
			mailGroup.POST("", middleware.Auth(), mail.SendMail)
			mailGroup.GET("/sent", middleware.Auth(), mail.RetrieveSentMails)
			mailGroup.POST("/delete", middleware.Auth(), mail.DeleteMails)
			mailGroup.POST("/send-administrator", middleware.AuthAdministrator(), mail.SendAdministratorMail)
			mailGroup.POST("/:id/set-read", middleware.Auth(), mail.SetMailRead)
		}

		{
			diaryGroup := api.Group("diaries")
			diaryGroup.GET("", middleware.Auth(), diary.RetrieveDiaryOverviews)
			diaryGroup.GET("/write", middleware.Auth(), diary.RetrieveDiaryEditData)
			diaryGroup.POST("/write", middleware.Auth(), diary.ReservePublishDiary)
			diaryGroup.POST("/write/clear", middleware.Auth(), diary.ClearReservedDiary)
			diaryGroup.GET("/write/preview", middleware.Auth(), diary.RetrieveDiaryPreview)
			diaryGroup.GET("/:nth/:character", middleware.Auth(), diary.RetrieveDiary)
		}

		{
			forumGroup := api.Group("forums")
			forumGroup.GET("", forum.RetrieveForumOverviews)
			forumGroup.GET("/:forum", forum.RetrieveForum)
			forumGroup.GET("/:forum/post-type", forum.RetrieveForumForcedPostType)
			forumGroup.GET("/:forum/topics", forum.RetrieveForumTopicOverviews)
			forumGroup.POST("/:forum/topics", forum.CreateForumTopic)
			forumGroup.GET("/:forum/topics/:topic", forum.RetrieveForumTopic)
			forumGroup.PUT("/:forum/topics/:topic", forum.UpdateForumTopic)
			forumGroup.POST("/:forum/topics/:topic/posts", forum.PostForumTopicPost)
			forumGroup.PUT("/:forum/topics/:topic/posts/:post", forum.ReviseForumTopicPost)
			forumGroup.POST("/:forum/topics/:topic/posts/:post/delete", forum.DeleteForumTopicPost)
			forumGroup.POST("/:forum/topics/:topic/posts/:post/reactions", middleware.Auth(), forum.ChangeForumReactionState)
		}

		{
			controlGroup := api.Group("control")

			{
				gameGroup := controlGroup.Group("game")
				gameGroup.POST("/announcements", middleware.AuthAdministrator(), control.Announce)
				gameGroup.GET("/announcements/:id", middleware.AuthAdministrator(), control.RetrieveAnnouncementEditData)
				gameGroup.POST("/announcements/:id/update", middleware.AuthAdministrator(), control.UpdateAnnouncement)
			}

			{
				debugGroup := controlGroup.Group("debug")
				debugGroup.POST("/dummy-character", middleware.AuthAdministrator(), control.CreateDummyCharacters)
			}
		}
	}

	return router
}
