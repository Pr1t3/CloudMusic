package models

import (
	"time"
)

type Album struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	AuthorId    int       `json:"author_id"`
	ReleaseDate time.Time `json:"release_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
