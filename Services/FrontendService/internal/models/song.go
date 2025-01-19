package models

import (
	"database/sql"
	"time"
)

type Song struct {
	Id       int           `json:"id"`
	Title    string        `json:"title"`
	Duration int           `json:"duration"`
	AlbumId  sql.NullInt32 `json:"album_id"`
	GenreId  sql.NullInt32 `json:"genre_id"`
	FilePath string        `json:"file_path"`
	AddedAt  time.Time     `json:"added_at"`
	Authors  []Author      `json:"authors"`
}
