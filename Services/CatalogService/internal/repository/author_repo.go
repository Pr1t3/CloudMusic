package repository

import (
	"CatalogService/internal/models"
	"database/sql"
)

type AuthorRepo struct {
	Db *sql.DB
}

func NewAuthorRepo(db *sql.DB) *AuthorRepo {
	return &AuthorRepo{Db: db}
}

func (a *AuthorRepo) GetAuthorById(authorId int) (*models.Author, error) {
	query := `SELECT * FROM authors WHERE id = ?`
	author := models.Author{}
	err := a.Db.QueryRow(query, authorId).Scan(&author.Id, &author.UserId, &author.Name)
	if err != nil {
		return nil, err
	}

	return &author, err
}

func (a *AuthorRepo) GetAuthorByUserId(userId int) (*models.Author, error) {
	query := `SELECT * FROM authors WHERE user_id = ?`
	author := models.Author{}
	err := a.Db.QueryRow(query, userId).Scan(&author.Id, &author.UserId, &author.Name)
	if err != nil {
		return nil, err
	}

	return &author, err
}

func (a *AuthorRepo) GetAuthorByName(name string) (*models.Author, error) {
	query := `SELECT * FROM authors WHERE name = ?`
	author := models.Author{}
	err := a.Db.QueryRow(query, name).Scan(&author.Id, &author.UserId, &author.Name)
	if err != nil {
		return nil, err
	}

	return &author, err
}

func (a *AuthorRepo) AddAuthor(userId int, name string) (int, error) {
	query := `INSERT INTO authors (user_id, name) VALUES(?,?)`
	res, err := a.Db.Exec(query, userId, name)
	if err != nil {
		return 0, err
	}
	authorId, err := res.LastInsertId()
	return int(authorId), err
}
