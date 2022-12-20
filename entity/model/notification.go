package model

import "time"

type Notification struct {
	Type      string    `json:"type"`
	Icon      *string   `json:"icon"`
	Message   string    `json:"message"`
	Detail    string    `json:"detail"`
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}
