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

func (sa *SongAuthorRepo) GetAllSongsByAuthorId(authorId int) ([]models.Song, error) {
	query := `SELECT s.* FROM songs AS s JOIN song_authors AS sa ON s.id = sa.song_id WHERE sa.author_id = ?`
	songs := []models.Song{}
	rows, err := sa.Db.Query(query, authorId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		song := models.Song{}
		if err := rows.Scan(&song.Id, &song.Title, &song.Duration, &song.Size, &song.GenreId, &song.FilePath, &song.AddedAt); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}
