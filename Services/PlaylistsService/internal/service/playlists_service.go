package service

import (
	"PlaylistsService/internal/models"
	"PlaylistsService/internal/repository"
)

type PlaylistService struct {
	playlistRepo        *repository.PlayListRepo
	songsInPlaylistRepo *repository.SongsInPlaylistRepo
}

func NewPlaylistService(p repository.PlayListRepo, sip repository.SongsInPlaylistRepo) *PlaylistService {
	return &PlaylistService{playlistRepo: &p, songsInPlaylistRepo: &sip}
}

func (ps *PlaylistService) GetPlaylistsByUserId(userId int) ([]models.Playlist, error) {
	return ps.playlistRepo.GetPlaylistsByUserId(userId)
}

func (ps *PlaylistService) GetPlaylistById(playlistId int) (*models.Playlist, error) {
	return ps.playlistRepo.GetPlaylistById(playlistId)
}

func (ps *PlaylistService) AddSongToPlaylist(playlistId, songId int) error {
	songs, err := ps.GetSongsInPlaylist(playlistId)
	if err != nil {
		return err
	}

	return ps.songsInPlaylistRepo.AddSongToPlaylist(playlistId, songId, len(songs)+1)
}

func (ps *PlaylistService) RemoveSongFromPlaylist(playlistId, songId int) error {
	return ps.songsInPlaylistRepo.RemoveSongFromPlaylist(playlistId, songId)
}

func (ps *PlaylistService) ChangePublicOption(playlistId int, isPublic bool) error {
	return ps.playlistRepo.ChangePublicOption(playlistId, isPublic)
}

func (ps *PlaylistService) GetPublicOption(playlistId int) (bool, error) {
	return ps.playlistRepo.GetPublicOption(playlistId)
}

func (ps *PlaylistService) GetSongsInPlaylist(playlistId int) ([]int, error) {
	return ps.songsInPlaylistRepo.GetSongsInPlaylist(playlistId)
}

func (ps *PlaylistService) RemovePlaylist(playlistId int) error {
	return ps.playlistRepo.RemovePlaylist(playlistId)
}

func (ps *PlaylistService) AddPlaylist(userId int, name string) (int, error) {
	return ps.playlistRepo.AddPlaylist(userId, name)
}
