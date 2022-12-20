package model

type DiaryOverview struct {
	Author CharacterOverview `json:"author"`
	Title  string            `json:"title"`
}

type Diary struct {
	DiaryOverview
	ExistingDiaries []int  `json:"existingDiaries"`
	Diary           string `json:"diary"`
	Nth             int    `json:"nth"`
}

type DiaryEditData struct {
	Title           string   `json:"title"`
	Diary           string   `json:"diary"`
	SelectableIcons []string `json:"selectableIcons"`
}
