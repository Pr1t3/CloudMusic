package handler

import (
	"CatalogService/internal/models"
	"CatalogService/internal/service"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hajimehoshi/go-mp3"
)

func getDuration(r []byte) (int, error) {
	reader := bytes.NewReader(r)

	d, err := mp3.NewDecoder(reader)
	if err != nil {
		return 0, err
	}
	const sampleSize = 4
	samples := int(d.Length() / sampleSize)
	audioLength := samples / d.SampleRate()
	return audioLength, nil
}

type CatalogHandler struct {
	catalogService *service.CatalogService
	*service.ProxyRequestStruct
}

func NewCatalogHandler(catalogService *service.CatalogService) *CatalogHandler {
	return &CatalogHandler{catalogService: catalogService}
}

func (h *CatalogHandler) GetSongs() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}
		songs, err := h.catalogService.GetSongs()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(songs)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
	})
}

func (h *CatalogHandler) GetSongById() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := r.URL.Path
		parts := strings.Split(path, "/")
		songId, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "Invalid song ID", http.StatusBadRequest)
			return
		}
		song, err := h.catalogService.GetSong(songId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(song)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
	})
}

func (h *CatalogHandler) GetAlbums() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		respBody, _, err := h.ProxyRequest(r, "http://localhost:9999/get-claims/", nil, http.MethodGet)
		if err != nil {
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}

		var claims models.Claims

		if err := json.Unmarshal(respBody, &claims); err != nil {
			http.Error(w, "Status Forbidden", http.StatusForbidden)
			return
		}

		albums, err := h.catalogService.GetAlbums(claims.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(albums)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
	})
}

func (h *CatalogHandler) GetGenres() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		genres, err := h.catalogService.GetGenres()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(genres)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
	})
}

func (h *CatalogHandler) GetSongAuthors() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}
		path := r.URL.Path
		parts := strings.Split(path, "/")
		songId, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "Invalid song ID", http.StatusBadRequest)
			return
		}
		authors, err := h.catalogService.GetAuthorsBySongId(songId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(authors)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
	})
}

func (h *CatalogHandler) AddSong() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		respBody, _, err := h.ProxyRequest(r, "http://localhost:9999/get-claims/", nil, http.MethodGet)
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
			http.Error(w, "Ошибка при чтении тела запроса", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		r.Body = io.NopCloser(io.Reader(bytes.NewReader(body)))

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Print(err.Error())
			http.Error(w, "Ошибка при парсинге формы", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		albumIDstring := r.FormValue("album_id")
		genreIDstring := r.FormValue("genre_id")
		var albumID *int
		if albumIDstring == "" {
			albumID = nil
		} else {
			res, err := strconv.Atoi(albumIDstring)
			if err != nil {
				log.Print(err.Error())
				http.Error(w, "Неверный ID альбома", http.StatusBadRequest)
				return
			}
			albumID = &res
		}

		genreID, err := strconv.Atoi(genreIDstring)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Неверный ID жанра", http.StatusBadRequest)
			return
		}

		authorsNames := r.Form["authors"]

		var authors []models.Author
		author, err := h.catalogService.GetAuthorByUserId(claims.UserId)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Не удалось найти автора", http.StatusInternalServerError)
			return
		}
		authors = append(authors, *author)
		for _, name := range authorsNames {
			trimmedName := strings.TrimSpace(name)
			author, err = h.catalogService.GetAuthorByName(trimmedName)
			if err != nil {
				log.Print(err.Error())
				http.Error(w, "Не удалось найти автора", http.StatusInternalServerError)
				return
			}
			authors = append(authors, *author)
		}

		file, _, err := r.FormFile("fileToUpload")
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
			return
		}
		defer file.Close()

		buf, err := io.ReadAll(file)
		if err != nil {
			log.Println(err)
			http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
			return
		}

		duration, err := getDuration(buf)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Не удалось извлечь длительность файла", http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(io.Reader(bytes.NewReader(body)))

		filePath := "/CloudMusic/" + strconv.Itoa(claims.UserId) + "/"

		r.Header.Add("FilePath", filePath)
		r.Header.Add("UserId", strconv.Itoa(claims.UserId))
		respBody, _, err = h.ProxyRequest(r, "http://localhost:9995/upload/", r.Body, http.MethodPost)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var responseData struct {
			FileType string
			FileName string
			Size     int64
		}

		if err := json.Unmarshal(respBody, &responseData); err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		songId, err := h.catalogService.AddSong(title, filePath+responseData.FileName, duration, albumID, &genreID)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, author := range authors {
			err = h.catalogService.AddAuthorBySongId(author, songId)
			if err != nil {
				log.Print(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})
}

func (h *CatalogHandler) BecomeAuthor() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		respBody, _, err := h.ProxyRequest(r, "http://localhost:9999/get-claims/", nil, http.MethodGet)
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

		name := r.FormValue("name")
		err = h.catalogService.AddAuthor(claims.UserId, name)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (h *CatalogHandler) IsAuthor() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		respBody, _, err := h.ProxyRequest(r, "http://localhost:9999/get-claims/", nil, http.MethodGet)
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

		_, err = h.catalogService.GetAuthorByUserId(claims.UserId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
