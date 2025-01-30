package handler

import (
	"SearchService/internal/service"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type SearchHandler struct {
	searchService *service.SearchService
}

func NewSearchHandler(ss service.SearchService) *SearchHandler {
	return &SearchHandler{searchService: &ss}
}

func (h *SearchHandler) InsertTerm() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		var bodyData struct {
			Term       string
			EntityId   int
			EntityType string
		}
		err = json.Unmarshal(body, &bodyData)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		err = h.searchService.InsertTerm(bodyData.Term, bodyData.EntityId, bodyData.EntityType)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (h *SearchHandler) SearchByPrefix() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		parts := strings.Split(r.URL.Path, "/")
		prefix := parts[len(parts)-1]
		entities, err := h.searchService.SearchPrefix(prefix)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonData, err := json.Marshal(entities)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})
}

func (h *SearchHandler) GetEntityByTerm() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		parts := strings.Split(r.URL.Path, "/")
		term := parts[len(parts)-1]
		entities, err := h.searchService.GetEntities(term)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonData, err := json.Marshal(entities)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})
}
