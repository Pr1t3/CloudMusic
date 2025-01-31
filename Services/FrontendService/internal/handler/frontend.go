package handler

import (
	"FrontendService/internal/models"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Claims struct {
	Email  string
	UserId int
	exp    time.Time
	iat    time.Time
}

func ProxyRequest(r *http.Request, target string, reqBody io.Reader, method string) ([]byte, *http.Header, error) {
	proxyReq, err := http.NewRequest(method, target, reqBody)
	if err != nil {
		return nil, nil, err
	}
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}
	for _, cookie := range r.Cookies() {
		proxyReq.AddCookie(cookie)
	}

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New("status forbidden")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, &resp.Header, err
}

func LoginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "./static/login.html")
	})
}

func NotAnAuthor() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "./static/notanauthorpage.html")
	})
}

func BecomeAuthorPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		_, _, err := ProxyRequest(r, "http://localhost:9989/is-author/", nil, http.MethodGet)
		if err == nil {
			http.Redirect(w, r, "http://localhost:9997/profile", http.StatusFound)
			return
		}
		log.Println("ASDASD")

		http.ServeFile(w, r, "./static/become_author_page.html")
	})
}

func RegisterHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		http.ServeFile(w, r, "./static/register.html")
	})
}

func ShowProfile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошибка при чтении тела запроса", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		r.Body = io.NopCloser(io.Reader(bytes.NewReader(body)))
		photoData, headers, err := ProxyRequest(r, "http://localhost:9999/get-profile-photo", nil, http.MethodGet)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		respBody, _, err := ProxyRequest(r, "http://localhost:9999/get-claims/", nil, http.MethodGet)

		if err != nil {
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}

		var claims Claims

		if err := json.Unmarshal(respBody, &claims); err != nil {
			http.Error(w, "Status Forbidden", http.StatusForbidden)
			return
		}

		templates := []string{
			"./static/profile.tmpl",
		}
		ts, err := template.ParseFiles(templates...)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var data struct {
			Email     string
			PhotoType string
			PhotoData string
			IsAuthor  bool
		}

		data.IsAuthor = true

		_, _, err = ProxyRequest(r, "http://localhost:9989/is-author/", nil, http.MethodGet)
		if err != nil {
			log.Println(err)
			data.IsAuthor = false
		}

		data.PhotoData = base64.StdEncoding.EncodeToString(photoData)
		data.PhotoType = headers.Get("Photo-Type")
		data.Email = claims.Email

		err = ts.Execute(w, data)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			return
		}
	})
}

func AddSong() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		templates := []string{
			"./static/add_song.tmpl",
		}
		ts, err := template.ParseFiles(templates...)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var data struct {
			Genres []models.Genre
		}

		respBody, _, err := ProxyRequest(r, "http://localhost:9989/genres", nil, http.MethodGet)
		if err != nil {
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(respBody, &data.Genres); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, data)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			return
		}
	})
}

func ShowMainPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		funcMap := template.FuncMap{
			"formatDuration": func(totalSeconds int) string {
				hours := totalSeconds / 3600
				minutes := totalSeconds / 60
				seconds := totalSeconds % 60
				if hours > 0 {
					return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
				}
				return fmt.Sprintf("%02d:%02d", minutes, seconds)
			},
		}

		templates := []string{
			"./static/index.tmpl",
		}
		ts, err := template.New("index.tmpl").Funcs(funcMap).ParseFiles(templates...)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		photoData, headers, err := ProxyRequest(r, "http://localhost:9999/get-profile-photo", r.Body, http.MethodGet)
		if err != nil {
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}

		var data struct {
			PhotoType string
			PhotoData string
			Songs     []models.Song
		}

		respBody, _, err := ProxyRequest(r, "http://localhost:9989/songs", nil, http.MethodGet)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(respBody, &data.Songs); err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		for index, song := range data.Songs {
			respBody, _, err = ProxyRequest(r, "http://localhost:9989/authors/"+fmt.Sprint(song.Id), nil, http.MethodGet)
			if err != nil {
				log.Print(err.Error())
				http.Error(w, "Server Internal Error", http.StatusInternalServerError)
				return
			}
			if err := json.Unmarshal(respBody, &data.Songs[index].Authors); err != nil {
				log.Print(err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		data.PhotoData = base64.StdEncoding.EncodeToString(photoData)
		data.PhotoType = headers.Get("Photo-Type")

		err = ts.Execute(w, data)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			return
		}
	})
}

func PlaylistsPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		templates := []string{
			"./static/all_playlists.tmpl",
		}

		funcMap := template.FuncMap{
			"FormatTime": func(t time.Time) string {
				return t.Format("January 2, 2006 at 15:04:05")
			},
		}

		ts, err := template.New("all_playlists.tmpl").Funcs(funcMap).ParseFiles(templates...)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var data struct {
			PhotoType string
			PhotoData string
			Playlists []models.Playlist
		}

		photoData, headers, err := ProxyRequest(r, "http://localhost:9999/get-profile-photo", nil, http.MethodGet)
		if err != nil {
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}

		data.PhotoData = base64.StdEncoding.EncodeToString(photoData)
		data.PhotoType = headers.Get("Photo-Type")

		respBody, _, err := ProxyRequest(r, "http://localhost:9987/playlists", nil, http.MethodGet)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(respBody, &data.Playlists); err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, data)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			return
		}
	})
}

func ShowPlaylist() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		templates := []string{
			"./static/playlist.tmpl",
		}

		funcMap := template.FuncMap{
			"formatDuration": func(totalSeconds int) string {
				hours := totalSeconds / 3600
				minutes := totalSeconds / 60
				seconds := totalSeconds % 60
				if hours > 0 {
					return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
				}
				return fmt.Sprintf("%02d:%02d", minutes, seconds)
			},
		}

		ts, err := template.New("playlist.tmpl").Funcs(funcMap).ParseFiles(templates...)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var data struct {
			PhotoType string
			PhotoData string
			Playlist  models.Playlist
			Songs     []models.Song
		}

		photoData, headers, err := ProxyRequest(r, "http://localhost:9999/get-profile-photo", nil, http.MethodGet)
		if err != nil {
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}

		data.PhotoData = base64.StdEncoding.EncodeToString(photoData)
		data.PhotoType = headers.Get("Photo-Type")

		parts := strings.Split(r.URL.Path, "/")
		playlistId := parts[len(parts)-1]

		respBody, _, err := ProxyRequest(r, "http://localhost:9987/playlists/"+playlistId, nil, http.MethodGet)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(respBody, &data); err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, data)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			return
		}
	})
}

func ShowAuthorPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		templates := []string{
			"./static/author.tmpl",
		}

		funcMap := template.FuncMap{
			"formatDuration": func(totalSeconds int) string {
				hours := totalSeconds / 3600
				minutes := totalSeconds / 60
				seconds := totalSeconds % 60
				if hours > 0 {
					return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
				}
				return fmt.Sprintf("%02d:%02d", minutes, seconds)
			},
		}

		ts, err := template.New("author.tmpl").Funcs(funcMap).ParseFiles(templates...)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var data struct {
			PhotoType string
			PhotoData string
			Author    models.Author
			Songs     []models.Song
		}

		photoData, headers, err := ProxyRequest(r, "http://localhost:9999/get-profile-photo", nil, http.MethodGet)
		if err != nil {
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}

		data.PhotoData = base64.StdEncoding.EncodeToString(photoData)
		data.PhotoType = headers.Get("Photo-Type")

		parts := strings.Split(r.URL.Path, "/")
		authorId := parts[len(parts)-1]
		respBody, _, err := ProxyRequest(r, "http://localhost:9989/songs-by-author/"+authorId, nil, http.MethodGet)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal(respBody, &data.Songs); err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		respBody, _, err = ProxyRequest(r, "http://localhost:9989/author/"+authorId, nil, http.MethodGet)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(respBody, &data.Author); err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			return
		}
	})
}
