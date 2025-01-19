package repository

import (
	"CatalogService/internal/models"
	"database/sql"
)

type SongAuthorRepo struct {
	Db *sql.DB
}

func NewSongAuthorRepo(db *sql.DB) *SongAuthorRepo {
	return &SongAuthorRepo{Db: db}
}

func (sa *SongAuthorRepo) GetAuthorsIdBySongId(songId int) ([]int, error) {
	query := `SELECT author_id FROM song_authors WHERE song_id = ?`
	rows, err := sa.Db.Query(query, songId)
	if err != nil {
		return nil, err
	}
	var authors []int
	for rows.Next() {
		var authorId int
		if err := rows.Scan(&authorId); err != nil {
			return nil, err
		}
		authors = append(authors, authorId)
	}
	return authors, nil
}

func (sa *SongAuthorRepo) AddAuthorBySongId(author models.Author, songId int) error {
	query := `INSERT INTO song_authors (song_id, author_id) VALUES(?, ?)`
	_, err := sa.Db.Exec(query, songId, author.Id)
	return err
}
