package repository

import (
	"CatalogService/internal/models"
	"database/sql"
)

type SongRepo struct {
	Db *sql.DB
}

func NewSongRepo(db *sql.DB) *SongRepo {
	return &SongRepo{Db: db}
}

func (s *SongRepo) GetSongs() ([]models.Song, error) {
	query := `SELECT * FROM songs ORDER BY added_at DESC LIMIT 10`
	rows, err := s.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var songs []models.Song
	for rows.Next() {
		song := models.Song{}
		if err := rows.Scan(&song.Id, &song.Title, &song.Duration, &song.AlbumId, &song.GenreId, &song.FilePath, &song.AddedAt); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}

func (s *SongRepo) GetSong(id int) (*models.Song, error) {
	query := `SELECT * FROM songs WHERE id = ?`
	row := s.Db.QueryRow(query, id)
	song := models.Song{}
	if err := row.Scan(&song.Id, &song.Title, &song.Duration, &song.AlbumId, &song.GenreId, &song.FilePath, &song.AddedAt); err != nil {
		return nil, err
	}
	return &song, nil
}

func (s *SongRepo) AddSong(title, filePath string, duration int, albumId, genreId *int) (int, error) {
	var query string
	var err error
	var songId int64
	var res sql.Result
	if albumId == nil {
		if genreId == nil {
			query = `INSERT INTO songs (title, duration, album_id, genre_id, file_path) VALUES(?, ?, null, null, ?)`
			res, err = s.Db.Exec(query, title, duration, filePath)
			if err != nil {
				return 0, err
			}
			songId, err = res.LastInsertId()
		} else {
			query = `INSERT INTO songs (title, duration, album_id, genre_id, file_path) VALUES(?, ?, null, ?, ?)`
			res, err = s.Db.Exec(query, title, duration, genreId, filePath)
			if err != nil {
				return 0, err
			}
			songId, err = res.LastInsertId()
		}
	} else {
		query = `INSERT INTO songs (title, duration, album_id, genre_id, file_path) VALUES(?, ?, ?, ?, ?)`
		res, err = s.Db.Exec(query, title, duration, albumId, genreId, filePath)
		if err != nil {
			return 0, err
		}
		songId, err = res.LastInsertId()
	}
	return int(songId), err
}
