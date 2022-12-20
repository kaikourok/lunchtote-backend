package model

import "time"

type AnnouncementBase struct {
	Id          int       `json:"id"`
	Type        string    `json:"type"`
	AnnouncedAt time.Time `json:"announcedAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type AnnouncementOverview struct {
	AnnouncementBase
	Overview string `json:"overview"`
}

type Announcement struct {
	AnnouncementBase
	Title   string `json:"title"`
	Content string `json:"content"`
}

type AnnouncementGuideData struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type AnnouncementEditData struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Overview string `json:"overview"`
	Content  string `json:"content"`
}

type AnnouncementEditDataUpdate struct {
	AnnouncementEditData
	SilentUpdate bool `json:"silentUpdate"`
}

type Inquiry struct {
	Id        int                `json:"id"`
	Character *CharacterOverview `json:"character"`
	Content   string             `json:"content"`
	Resolved  bool               `json:"resolved"`
	PostedAt  time.Time          `json:"postedAt"`
}

type AdomonishData struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

type BanData struct {
	Reason string `json:"reason"`
}

type UnbanData struct {
	Reason string `json:"reason"`
}

type ProhibitionRelatedData struct {
	Id        int       `json:"id"`
	Type      string    `json:"type"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}
