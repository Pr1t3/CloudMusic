package repository

import "database/sql"

type SongsInPlaylistRepo struct {
	Db *sql.DB
}

func NewSongsInPlaylistRepo(db *sql.DB) *SongsInPlaylistRepo {
	return &SongsInPlaylistRepo{Db: db}
}

func (sip *SongsInPlaylistRepo) AddSongToPlaylist(playlistId, songId, songOrder int) error {
	query := `INSERT INTO songs_in_playlists (playlist_id, song_id, song_order) VALUES(?,?,?)`
	_, err := sip.Db.Exec(query, playlistId, songId, songOrder)
	return err
}

func (sip *SongsInPlaylistRepo) RemoveSongFromPlaylist(playlistId, songId int) error {
	query := `DELETE FROM songs_in_playlists WHERE playlist_id =? AND song_id =?`
	_, err := sip.Db.Exec(query, playlistId, songId)
	return err
}

func (sip *SongsInPlaylistRepo) GetSongsInPlaylist(playlistId int) ([]int, error) {
	query := `SELECT song_id FROM songs_in_playlists WHERE playlist_id = ? ORDER BY song_order ASC`
	rows, err := sip.Db.Query(query, playlistId)
	if err != nil {
		return nil, err
	}
	var songIds []int
	for rows.Next() {
		var songId int
		if err := rows.Scan(&songId); err != nil {
			return nil, err
		}
		songIds = append(songIds, songId)
	}
	return songIds, nil
}
