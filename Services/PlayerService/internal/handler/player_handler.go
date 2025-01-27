package handler

import (
	"PlayerService/internal/models"
	"PlayerService/internal/service"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func StartSong() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		songId, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "Invalid song ID", http.StatusBadRequest)
			return
		}

		var reqBodyStruct struct {
			FilePath string `json:"filePath"`
		}

		respBody, _, err := service.ProxyRequest(r, "http://localhost:9989/songs/"+strconv.Itoa(songId), nil, http.MethodGet)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error starting song", http.StatusInternalServerError)
			return
		}

		var song models.Song
		err = json.Unmarshal(respBody, &song)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error starting song", http.StatusInternalServerError)
			return
		}

		reqBodyStruct.FilePath = song.FilePath
		reqBody, err := json.Marshal(reqBodyStruct)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error starting song", http.StatusInternalServerError)
			return
		}

		respBody, respHeader, err := service.ProxyRequest(r, "http://localhost:9995/download", bytes.NewBuffer(reqBody), http.MethodGet)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error starting song", http.StatusInternalServerError)
			return
		}

		for key, values := range *respHeader {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(http.StatusPartialContent)
		w.Write(respBody)
	})
}
