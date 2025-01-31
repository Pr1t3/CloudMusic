package handler

import (
	"PlaylistsService/internal/models"
	"PlaylistsService/internal/service"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type PlaylistHandler struct {
	playlistService *service.PlaylistService
	*service.ProxyRequestStruct
}

func NewPlaylistHandler(ps *service.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{playlistService: ps}
}

func (p *PlaylistHandler) CreatePlaylist() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		respBody, _, err := p.ProxyRequest(r, "http://localhost:9999/get-claims/", nil, http.MethodGet)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		var claims models.Claims

		if err := json.Unmarshal(respBody, &claims); err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var playlistName struct {
			Name string
		}

		err = json.Unmarshal(body, &playlistName)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Ошибка при распаковке тела запроса", http.StatusBadRequest)
			return
		}

		playlistId, err := p.playlistService.AddPlaylist(claims.UserId, playlistName.Name)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Ошибка при добавлении плейлиста", http.StatusInternalServerError)
			return
		}

		var reqBody struct {
			Term       string
			EntityId   int
			EntityType string
		}
		reqBody.Term = playlistName.Name
		reqBody.EntityId = playlistId
		reqBody.EntityType = "playlist"
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, _, err = p.ProxyRequest(r, "http://localhost:9986/add-term/", bytes.NewBuffer(jsonData), http.MethodPost)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func (p *PlaylistHandler) DeletePlaylist() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		playlistId, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
			return
		}

		err = p.playlistService.RemovePlaylist(playlistId)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Ошибка при удалении плейлиста", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (p *PlaylistHandler) ChangePublicOption(isPublic bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		playlistId, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
			return
		}

		err = p.playlistService.ChangePublicOption(playlistId, isPublic)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Ошибка при изменении публичности плейлиста", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (p *PlaylistHandler) GetPlaylists() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		respBody, _, err := p.ProxyRequest(r, "http://localhost:9999/get-claims/", nil, http.MethodGet)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		var claims models.Claims

		if err := json.Unmarshal(respBody, &claims); err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		playlists, err := p.playlistService.GetPlaylistsByUserId(claims.UserId)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Ошибка при получении плейлистов", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(playlists)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
	})
}

func (p *PlaylistHandler) GetPlaylistById() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		playlistId, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
			return
		}

		songsIds, err := p.playlistService.GetSongsInPlaylist(playlistId)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Ошибка при получении песен из плейлиста", http.StatusInternalServerError)
			return
		}

		curPlaylist, err := p.playlistService.GetPlaylistById(playlistId)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Ошибка при получении плейлиста", http.StatusInternalServerError)
			return
		}

		var data struct {
			Playlist models.Playlist
			Songs    []models.Song
		}

		data.Playlist = *curPlaylist

		for _, songId := range songsIds {

			respBody, _, err := p.ProxyRequest(r, "http://localhost:9989/songs/"+strconv.Itoa(songId), nil, http.MethodGet)
			if err != nil {
				log.Print(err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			var song models.Song
			if err := json.Unmarshal(respBody, &song); err != nil {
				log.Print(err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			data.Songs = append(data.Songs, song)
		}

		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
	})
}

func (p *PlaylistHandler) AddSongToPlaylist() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		playlistId, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var songId struct {
			SongId int
		}

		err = json.Unmarshal(body, &songId)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Ошибка при распаковке тела запроса", http.StatusBadRequest)
			return
		}

		err = p.playlistService.AddSongToPlaylist(playlistId, songId.SongId)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Ошибка при добавлении песни в плейлист", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (p *PlaylistHandler) RemoveSongFromPlaylist() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		playlistId, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var songId struct {
			SongId int
		}

		err = json.Unmarshal(body, &songId)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Ошибка при распаковке тела запроса", http.StatusBadRequest)
			return
		}

		err = p.playlistService.RemoveSongFromPlaylist(playlistId, songId.SongId)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Ошибка при удалении песни из плейлиста", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
