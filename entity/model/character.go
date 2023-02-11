package model

import (
	"time"

	"github.com/kaikourok/lunchtote-backend/entity/service"
)

type CharacterOverview struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Mainicon string `json:"mainicon"`
}

type Icon struct {
	Path string `json:"path"`
}

type ProfileImage struct {
	Path string `json:"path"`
}

type ProfileEditData struct {
	Name      string   `json:"name"`
	Nickname  string   `json:"nickname"`
	Summary   string   `json:"summary"`
	Profile   string   `json:"profile"`
	Mainicon  string   `json:"mainicon"`
	ListImage string   `json:"listImage"`
	Tags      []string `json:"tags"`
}

type ProfileDiaryData struct {
	Nth   int    `json:"nth"`
	Title string `json:"title"`
}

type Profile struct {
	Id              int                `json:"id"`
	Name            string             `json:"name"`
	Nickname        string             `json:"nickname"`
	Summary         string             `json:"summary"`
	Profile         string             `json:"profile"`
	ProfileImages   []string           `json:"profileImages"`
	Tags            []string           `json:"tags"`
	FollowingNumber int                `json:"followingNumber"`
	FollowedNumber  int                `json:"followedNumber"`
	Icons           []Icon             `json:"icons"`
	IsFollowing     bool               `json:"isFollowing"`
	IsFollowed      bool               `json:"isFollowed"`
	IsMuting        bool               `json:"isMuting"`
	IsBlocking      bool               `json:"isBlocking"`
	IsBlocked       bool               `json:"isBlocked"`
	ExistingDiaries []ProfileDiaryData `json:"existingDiaries"`
}

type AllCharacterListItem struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Nickname    string   `json:"nickname"`
	Summary     string   `json:"summary"`
	ListImage   string   `json:"listImage"`
	Tags        []string `json:"tags"`
	IsFollowing *bool    `json:"isFollowing,omitempty"`
	IsFollowed  *bool    `json:"isFollowed,omitempty"`
	IsMuting    *bool    `json:"isMuting,omitempty"`
	IsBlocking  *bool    `json:"isBlocking,omitempty"`
}

type GeneralCharacterListItem struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Nickname    string   `json:"nickname"`
	Summary     string   `json:"summary"`
	Mainicon    string   `json:"mainicon"`
	Tags        []string `json:"tags"`
	IsFollowing *bool    `json:"isFollowing,omitempty"`
	IsFollowed  *bool    `json:"isFollowed,omitempty"`
	IsMuting    *bool    `json:"isMuting,omitempty"`
	IsBlocking  *bool    `json:"isBlocking,omitempty"`
}

type UploadedImage struct {
	Id         int       `json:"id"`
	Path       string    `json:"path"`
	UploadedAt time.Time `json:"uploadedAt"`
}

type HomeData struct {
	Nickname string `json:"nickname"`
}

type CharacterSuggestionData struct {
	Id   int
	Name string
}

type CharacterSuggestion struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

type CharacterSuggestionsData []CharacterSuggestionData
type CharacterSuggestions []CharacterSuggestion

func (s *CharacterSuggestionData) ToDomain() *CharacterSuggestion {
	return &CharacterSuggestion{
		Id:   s.Id,
		Text: service.ConvertCharacterIdToText(s.Id) + " " + s.Name,
	}
}

func (s *CharacterSuggestionsData) ToDomain() *CharacterSuggestions {
	suggestions := make(CharacterSuggestions, len(*s))
	for i, suggestionData := range *s {
		suggestions[i] = *suggestionData.ToDomain()
	}
	return &suggestions
}

type CharacterEmailRegistratedData struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type CharacterIconLayerItemEditData struct {
	Path string `json:"path"`
}

type CharacterIconLayerItem struct {
	CharacterIconLayerItemEditData
	Id int `json:"id"`
}

type CharacterIconProcessSchema struct {
	Name   string  `json:"name"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Rotate float64 `json:"rotate"` // degree
	Scale  float64 `json:"scale"`  // percent
}

type CharacterIconProcessSchemaEditData struct {
	CharacterIconProcessSchema
	Id int `json:"id"`
}

type CharacterIconLayerGroup struct {
	Id    int                      `json:"id"`
	Name  string                   `json:"name"`
	Items []CharacterIconLayerItem `json:"items"`
}

type CharacterIconLayeringGroupOverview struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CharacterIconLayeringGroup struct {
	CharacterIconLayeringGroupOverview
	Layers    []CharacterIconLayerGroup            `json:"layers"`
	Processes []CharacterIconProcessSchemaEditData `json:"processes"`
}

type CharacterIconLayerGroupOrder struct {
	LayerGroup int `json:"layerGroup"`
	Order      int `json:"order"`
}

type CharacterOtherSettings struct {
	Webhook struct {
		Url       string `json:"url"`
		Followed  bool   `json:"followed"`
		Replied   bool   `json:"replied"`
		Subscribe bool   `json:"subscribe"`
		NewMember bool   `json:"newMember"`
		Mail      bool   `json:"mail"`
	} `json:"webhook"`
	Notification struct {
		Followed  bool `json:"followed"`
		Replied   bool `json:"replied"`
		Subscribe bool `json:"subscribe"`
		NewMember bool `json:"newMember"`
		Mail      bool `json:"mail"`
	} `json:"notification"`
}

type CharacterOtherSettingsState struct {
	Email        *string `json:"email"`
	LinkedStates struct {
		Twitter bool `json:"twitter"`
		Google  bool `json:"google"`
	} `json:"linkedStates"`
	CharacterOtherSettings
}
