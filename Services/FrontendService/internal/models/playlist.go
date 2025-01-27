package models

import "time"

type Playlist struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	UserId    int       `json:"user_id"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
}
