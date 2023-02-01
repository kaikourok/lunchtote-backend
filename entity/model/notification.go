package model

import "time"

const (
	NotificationTypeFollowed  = "FOLLOWED"
	NotificationTypeReplied   = "REPLIED"
	NotificationTypeSubscribe = "SUBSCRIBE"
	NotificationTypeNewMember = "NEW_MEMBER"
	NotificationTypeMail      = "MAIL"
	NotificationTypeMassMail  = "MASS_MAIL"
)

type Notification interface {
	implementsNotificationBase()
}

type NotificationBase struct {
	Id        int       `json:"id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

func (n NotificationBase) implementsNotificationBase() {}

type FollowedNotificationValue struct {
	Character CharacterOverview `json:"character"`
}

type FollowedNotification struct {
	NotificationBase
	FollowedNotificationValue
}

type RepliedNotificationValue struct {
	Character CharacterOverview `json:"character"`
	Room      RoomOverview      `json:"room"`
	Message   struct {
		ReferRoot int    `json:"referRoot"`
		Message   string `json:"message"`
	} `json:"message"`
}

type RepliedNotification struct {
	NotificationBase
	RepliedNotificationValue
}

type SubscribeNotificationValue struct {
	Character CharacterOverview `json:"character"`
	Room      RoomOverview      `json:"room"`
	Message   struct {
		Message string `json:"message"`
	} `json:"message"`
}

type SubscribeNotification struct {
	NotificationBase
	SubscribeNotificationValue
}

type NewMemberNotificationValue struct {
	Character CharacterOverview `json:"character"`
	Room      RoomOverview      `json:"room"`
}

type NewMemberNotification struct {
	NotificationBase
	NewMemberNotificationValue
}

type MailNotificationValue struct {
	Character struct {
		Id       *int    `json:"id"`
		Name     string  `json:"name"`
		Mainicon *string `json:"mainicon"`
	} `json:"character"`
	Mail struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
	} `json:"mail"`
}

type MailNotification struct {
	NotificationBase
	MailNotificationValue
}

type MassMailNotificationValue struct {
	Mail struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Title string `json:"title"`
	} `json:"mail"`
}

type MassMailNotification struct {
	NotificationBase
	MassMailNotificationValue
}
