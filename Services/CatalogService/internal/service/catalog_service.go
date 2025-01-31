package service

import (
	"CatalogService/internal/models"
	"CatalogService/internal/repository"
)

type CatalogService struct {
	songRepo       *repository.SongRepo
	authorRepo     *repository.AuthorRepo
	genreRepo      *repository.GenreRepo
	songAuthorRepo *repository.SongAuthorRepo
}

func NewCatalogService(s repository.SongRepo, au repository.AuthorRepo, g repository.GenreRepo, sa repository.SongAuthorRepo) *CatalogService {
	return &CatalogService{songRepo: &s, authorRepo: &au, genreRepo: &g, songAuthorRepo: &sa}
}

func NewCatalogServiceOnlyAuthor(au repository.AuthorRepo) *CatalogService {
	return &CatalogService{authorRepo: &au}
}

func (c *CatalogService) GetSongs() ([]models.Song, error) {
	return c.songRepo.GetSongs()
}

func (c *CatalogService) GetSong(id int) (*models.Song, error) {
	return c.songRepo.GetSong(id)
}

func (c *CatalogService) AddSong(title, filePath string, duration int, genreId *int, size int64) (int, error) {
	return c.songRepo.AddSong(title, filePath, duration, genreId, size)
}

func (c *CatalogService) AddAuthorBySongId(author models.Author, songId int) error {
	return c.songAuthorRepo.AddAuthorBySongId(author, songId)
}

func (c *CatalogService) GetAuthorByName(name string) (*models.Author, error) {
	return c.authorRepo.GetAuthorByName(name)
}

func (c *CatalogService) GetAuthorByUserId(userId int) (*models.Author, error) {
	return c.authorRepo.GetAuthorByUserId(userId)
}

func (c *CatalogService) GetGenres() ([]models.Genre, error) {
	return c.genreRepo.GetGenres()
}

func (c *CatalogService) GetAuthorsBySongId(songId int) ([]models.Author, error) {
	authorsId, err := c.songAuthorRepo.GetAuthorsIdBySongId(songId)
	if err != nil {
		return nil, err
	}

	var authors []models.Author
	for _, id := range authorsId {
		author, err := c.authorRepo.GetAuthorById(id)
		if err != nil {
			return nil, err
		}
		authors = append(authors, *author)
	}

	return authors, nil
}

func (c *CatalogService) GetAuthorById(id int) (*models.Author, error) {
	return c.authorRepo.GetAuthorById(id)
}

func (c *CatalogService) AddAuthor(userId int, name string) (int, error) {
	return c.authorRepo.AddAuthor(userId, name)
}

func (c *CatalogService) GetAllSongsByAuthorId(authorId int) ([]models.Song, error) {
	return c.songAuthorRepo.GetAllSongsByAuthorId(authorId)
}
