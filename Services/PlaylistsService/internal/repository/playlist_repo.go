package repository

import (
	"PlaylistsService/internal/models"
	"database/sql"
)

type PlayListRepo struct {
	Db *sql.DB
}

func NewPlaylistRepo(db *sql.DB) *PlayListRepo {
	return &PlayListRepo{Db: db}
}

func (p *PlayListRepo) GetPlaylistsByUserId(userId int) ([]models.Playlist, error) {
	query := `SELECT * FROM playlists WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := p.Db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	var playlists []models.Playlist
	for rows.Next() {
		playlist := models.Playlist{}
		if err := rows.Scan(&playlist.Id, &playlist.Name, &playlist.UserId, &playlist.IsPublic, &playlist.CreatedAt); err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}

func (p *PlayListRepo) GetPlaylistById(playlistId int) (*models.Playlist, error) {
	query := `SELECT * FROM playlists WHERE id = ?`
	row := p.Db.QueryRow(query, playlistId)
	playlist := models.Playlist{}
	if err := row.Scan(&playlist.Id, &playlist.Name, &playlist.UserId, &playlist.IsPublic, &playlist.CreatedAt); err != nil {
		return nil, err
	}
	return &playlist, nil
}

func (p *PlayListRepo) AddPlaylist(userId int, name string) (int, error) {
	query := `INSERT INTO playlists (user_id, name) VALUES(?,?)`
	res, err := p.Db.Exec(query, userId, name)
	if err != nil {
		return 0, err
	}
	playlistId, err := res.LastInsertId()
	return int(playlistId), err
}

func (p *PlayListRepo) RemovePlaylist(playlistId int) error {
	query := `DELETE FROM playlists WHERE id =?`
	_, err := p.Db.Exec(query, playlistId)
	return err
}

func (p *PlayListRepo) ChangePublicOption(playlistId int, isPublic bool) error {
	query := `UPDATE playlists SET is_public = ? WHERE id = ?`
	_, err := p.Db.Exec(query, isPublic, playlistId)
	return err
}

func (p *PlayListRepo) GetPublicOption(playlistId int) (bool, error) {
	query := `SELECT is_public FROM playlists WHERE id = ?`
	var isPublic bool
	err := p.Db.QueryRow(query, playlistId).Scan(&isPublic)
	if err != nil {
		return false, err
	}
	return isPublic, nil
}
