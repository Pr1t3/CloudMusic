package models

type SongsInPlaylist struct {
	PlaylistId int `json:"playlist_id"`
	SongId     int `json:"song_id"`
	SongOrder  int `json:"song_order"`
}
