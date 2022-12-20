package model

import "time"

type ForumOverviewLastPost struct {
	Topic struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
	} `json:"topic"`
	Sender struct {
		PostType  string `json:"postType"`
		Name      string `json:"name"`
		Character *int   `json:"character"`
	} `json:"sender"`
	PostedAt time.Time `json:"postedAt"`
}

type ForumOverview struct {
	Id       int                    `json:"id"`
	Title    string                 `json:"title"`
	Summary  string                 `json:"summary"`
	LastPost *ForumOverviewLastPost `json:"lastPost"`
}

type ForumGroup struct {
	Id     int             `json:"id"`
	Title  string          `json:"title"`
	Forums []ForumOverview `json:"forums"`
}

type Forum struct {
	Id             int     `json:"id"`
	Title          string  `json:"title"`
	Summary        string  `json:"summary"`
	Guide          string  `json:"guide"`
	ForcedPostType *string `json:"forcedPostType"`
}

type ForumCreateData struct {
	Title         string  `json:"title"`
	Summary       string  `json:"summary"`
	Guide         string  `json:"guide"`
	ForcePostType *string `json:"forcePostType"`
}

type ForumSender struct {
	PostType   string  `json:"postType"`
	Name       *string `json:"name"`
	Character  *int    `json:"character"`
	Identifier *string `json:"identifier"`
}

type ForumTopicOverview struct {
	Id           int         `json:"id"`
	Sender       ForumSender `json:"sender"`
	Title        string      `json:"title"`
	Status       string      `json:"status"`
	Posts        int         `json:"posts"`
	CreatedAt    time.Time   `json:"createdAt"`
	LastPostedAt time.Time   `json:"lastPostedAt"`
	IsPinned     bool        `json:"isPinned"`
}

type ForumTopic struct {
	Id            int         `json:"id"`
	Forum         int         `json:"forum"`
	Sender        ForumSender `json:"sender"`
	Title         string      `json:"title"`
	Status        string      `json:"status"`
	ForcePostType *string     `json:"forcePostType"`
}

type ForumTopicEditData struct {
	Title        string  `json:"title"`
	Status       string  `json:"status"`
	PostType     string  `json:"postType"`
	EditPassword *string `json:"editPassword"`
}

type ForumTopicCreateData struct {
	Title        string  `json:"title"`
	Name         *string `json:"name"`
	Content      string  `json:"content"`
	EditPassword *string `json:"editPassword"`
	PostType     string  `json:"postType"`
}

type ForumTopicPostRevisionHistory struct {
	Content   string    `json:"content"`
	RevisedAt time.Time `json:"revisedAt"`
}

type ForumTopicPostReaction struct {
	Emoji         string `json:"emoji"`
	ReactedCounts int    `json:"reactedCounts"`
	IsReacted     bool   `json:"isReacted"`
}

type ForumTopicPostRevision struct {
	Content  string    `json:"content"`
	PostedAt time.Time `json:"postedAt"`
}

type ForumTopicPost struct {
	Id        int                      `json:"id"`
	Sender    ForumSender              `json:"sender"`
	Content   string                   `json:"content"`
	PostedAt  time.Time                `json:"postedAt"`
	UpdatedAt *time.Time               `json:"updatedAt"`
	Revisions []ForumTopicPostRevision `json:"revisions"`
	Reactions []ForumTopicPostReaction `json:"reactions"`
}

type ForumTopicPostSendData struct {
	PostType     string  `json:"postType"`
	Name         *string `json:"name"`
	Content      string  `json:"content"`
	EditPassword *string `json:"editPassword"`
}

type ForumTopicPostReviseData struct {
	Content      string  `json:"content"`
	PostType     string  `json:"postType"`
	EditPassword *string `json:"editPassword"`
}
