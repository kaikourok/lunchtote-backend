package model

import "time"

type Room struct {
	Title               string   `json:"title"`
	Summary             string   `json:"summary"`
	Description         string   `json:"description"`
	Tags                []string `json:"tags"`
	Searchable          bool     `json:"searchable"`
	AllowRecommendation bool     `json:"allowRecommendation"`
	ChildrenReferable   bool     `json:"childrenReferable"`
	ParentRoom          *int     `json:"parentRoom"`
}

type RoomDetailData struct {
	Room
	Official     bool              `json:"official"`
	UpdatedAt    time.Time         `json:"updatedAt"`
	MembersCount int               `json:"membersCount"`
	Master       CharacterOverview `json:"master"`
}

type RoomOverview struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type RoomRelations struct {
	Parent   *RoomOverview  `json:"parent"`
	Siblings []RoomOverview `json:"siblings"`
	Children []RoomOverview `json:"children"`
}

type RoomListItem struct {
	Id     int `json:"id"`
	Master struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"master"`
	Title           string              `json:"title"`
	Summary         string              `json:"summary"`
	Tags            []string            `json:"tags"`
	Official        bool                `json:"official"`
	MessagesCount   int                 `json:"messagesCount"`
	MembersCount    int                 `json:"membersCount"`
	LastUpdate      time.Time           `json:"lastUpdate"`
	PostsPerDay     float64             `json:"postsPerDay"`
	FollowedMembers []CharacterOverview `json:"followedMembers,omitempty"`
}

type RoomRolePermission struct {
	Write              *bool `json:"write"`
	Ban                *bool `json:"ban"`
	Invite             *bool `json:"invite"`
	UseReply           *bool `json:"useReply"`
	UseSecret          *bool `json:"useSecret"`
	DeleteOtherMessage *bool `json:"deleteOtherMessage"`
	CreateChildrenRoom *bool `json:"createChildrenRoom"`
}

type RoomMemberPermission struct {
	Write              bool `json:"write"`
	Ban                bool `json:"ban"`
	Invite             bool `json:"invite"`
	UseReply           bool `json:"useReply"`
	UseSecret          bool `json:"useSecret"`
	DeleteOtherMessage bool `json:"deleteOtherMessage"`
	CreateChildrenRoom bool `json:"createChildrenRoom"`
}

type RoomRole struct {
	Id       int    `json:"id"`
	Priority int    `json:"priority"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	RoomRolePermission
}

type RoomRolePriority struct {
	Role     int `json:"Role"`
	Priority int `json:"Priority"`
}

type RoomRoleOverview struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type RomeMemberWithRoles struct {
	CharacterOverview
	Roles []RoomRoleOverview `json:"roles"`
}

type RoomInviteState struct {
	Invited   CharacterOverview `json:"invited"`
	Inviter   CharacterOverview `json:"inviter"`
	InvitedAt time.Time         `json:"invitedAt"`
}

type RoomBanState struct {
	Banned   CharacterOverview `json:"banned"`
	Banner   CharacterOverview `json:"banner"`
	BannedAt time.Time         `json:"bannedAt"`
}

type RoomMessage struct {
	Id              int          `json:"id"`
	Character       int          `json:"character"`
	Refer           *int         `json:"refer"`
	ReferRoot       *int         `json:"referRoot"`
	Secret          bool         `json:"secret"`
	Icon            *string      `json:"icon"`
	Name            string       `json:"name"`
	Message         string       `json:"message"`
	RepliedCount    int          `json:"repliedCount"`
	PostedAt        time.Time    `json:"postedAt"`
	ReplyPermission string       `json:"replyPermission"`
	Replyable       bool         `json:"replyable"`
	Room            RoomOverview `json:"room"`
	Recipients      []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"recipients"`
}

type RoomMessageRetrieveSettings struct {
	RangeType         string
	BasePoint         int
	FetchNumber       int
	Category          string
	Room              *int
	ReferRoot         *int
	Search            *string
	ListId            *int
	TargetCharacterId *int
	RelateFilter      bool
	Children          bool
}

type RoomPostMessage struct {
	Room            int    `json:"room"`
	Icon            string `json:"icon"`
	Name            string `json:"name"`
	Message         string `json:"message"`
	Refer           *int   `json:"refer"`
	DirectReply     *int   `json:"directReply"`
	ReplyPermission string `json:"replyPermission"`
	Secret          bool   `json:"secret"`
}

type RoomMessageEditRequiredData struct {
	Character struct {
		Name string `json:"name"`
	} `json:"character"`
	Icons        []Icon                         `json:"icons"`
	Lists        []ListOverview                 `json:"lists"`
	FetchConfigs []RoomMessageFetchConfigWithId `json:"fetchConfigs"`
}

type RoomMessageFetchConfig struct {
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	Room         *int    `json:"room"`
	Search       *string `json:"search"`
	ReferRoot    *int    `json:"referRoot"`
	List         *int    `json:"list"`
	Character    *int    `json:"character"`
	RelateFilter *bool   `json:"relateFilter"`
	Children     *bool   `json:"children"`
}

type RoomMessageFetchConfigWithId struct {
	Id int `json:"id"`
	RoomMessageFetchConfig
}

type RoomMessageFetchConfigOrder struct {
	Config int `json:"config"`
	Order  int `json:"order"`
}

type RoomSearchOptions struct {
	Title        string   `json:"title"`
	Tags         []string `json:"tags"`
	ExcludedTags []string `json:"excludedTags"`
	Description  string   `json:"description"`
	Order        string   `json:"order"`
	Sort         string   `json:"sort"`
	Participant  *string  `json:"participant"`
	Offset       int      `json:"offset"`
	Limit        int      `json:"limit"`
}
