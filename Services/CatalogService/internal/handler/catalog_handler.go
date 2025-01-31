package handler

import (
	"CatalogService/internal/models"
	"CatalogService/internal/service"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Eyevinn/mp4ff/mp4"
)

func getDuration(buf []byte) (int, error) {
	reader := bytes.NewReader(buf)

	mp4File, err := mp4.DecodeFile(reader)
	if err != nil {
		return 0, errors.New("ошибка при декодировании MP4 файла: " + err.Error())
	}

	if mp4File.Moov == nil || len(mp4File.Moov.Traks) == 0 {
		return 0, errors.New("недостаточно метаданных в MP4 файле")
	}

	track := mp4File.Moov.Traks[0]
	mdhd := track.Mdia.Mdhd
	if mdhd == nil {
		return 0, errors.New("отсутствуют метаданные MDHD")
	}

	duration := mdhd.Duration
	timeScale := mdhd.Timescale

	return int(duration / uint64(timeScale)), nil
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
		if r.Method != http.MethodGet {
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
		genreIDstring := r.FormValue("genre_id")

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

		songId, err := h.catalogService.AddSong(title, filePath+responseData.FileName, duration, &genreID, responseData.Size)
		if err != nil {
			log.Print(err)
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

		var reqBody struct {
			Term       string
			EntityId   int
			EntityType string
		}
		reqBody.Term = title
		reqBody.EntityId = songId
		reqBody.EntityType = "song"
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, _, err = h.ProxyRequest(r, "http://localhost:9986/add-term/", bytes.NewBuffer(jsonData), http.MethodPost)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
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
		authorId, err := h.catalogService.AddAuthor(claims.UserId, name)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var reqBody struct {
			Term       string
			EntityId   int
			EntityType string
		}
		reqBody.Term = name
		reqBody.EntityId = authorId
		reqBody.EntityType = "author"
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, _, err = h.ProxyRequest(r, "http://localhost:9986/add-term/", bytes.NewBuffer(jsonData), http.MethodPost)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

func (h *CatalogHandler) GetAuthorInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		authorId, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Invalid author ID", http.StatusBadRequest)
			return
		}

		author, err := h.catalogService.GetAuthorById(authorId)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(author)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
	})
}

func (h *CatalogHandler) GetAllSongsByAuthor() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		authorId, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Invalid author ID", http.StatusBadRequest)
			return
		}

		songs, err := h.catalogService.GetAllSongsByAuthorId(authorId)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonData, err := json.Marshal(songs)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})
}
