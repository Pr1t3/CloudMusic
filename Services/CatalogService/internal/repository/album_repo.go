package repository

import (
	"CatalogService/internal/models"
	"database/sql"
)

type AlbumRepo struct {
	Db *sql.DB
}

func NewAlbumRepo(db *sql.DB) *AlbumRepo {
	return &AlbumRepo{Db: db}
}

func (a *AlbumRepo) GetAlbumsByAuthorId(authorId int) ([]models.Album, error) {
	query := `SELECT * FROM albums WHERE author_id = ?`
	rows, err := a.Db.Query(query, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var albums []models.Album
	for rows.Next() {
		album := models.Album{}
		if err := rows.Scan(&album.Id, &album.Title, &album.AuthorId, &album.ReleaseDate, &album.CreatedAt, &album.UpdatedAt); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}
	return albums, nil
}
