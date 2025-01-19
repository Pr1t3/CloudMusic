package repository

import (
	"CatalogService/internal/models"
	"database/sql"
)

type GenreRepo struct {
	Db *sql.DB
}

func NewGenreRepo(db *sql.DB) *GenreRepo {
	return &GenreRepo{Db: db}
}

func (a *GenreRepo) GetGenres() ([]models.Genre, error) {
	query := `SELECT * FROM genres`
	rows, err := a.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var Genres []models.Genre
	for rows.Next() {
		Genre := models.Genre{}
		if err := rows.Scan(&Genre.Id, &Genre.Name); err != nil {
			return nil, err
		}
		Genres = append(Genres, Genre)
	}
	return Genres, nil
}
