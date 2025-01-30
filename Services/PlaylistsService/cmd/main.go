package main

import (
	"PlaylistsService/internal/handler"
	"PlaylistsService/internal/repository"
	"PlaylistsService/internal/service"
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()
	db, err := openDB(os.Args[1])
	if err != nil {
		print("Error in opening db")
	}
	defer db.Close()

	playlistHandler := handler.NewPlaylistHandler(service.NewPlaylistService(*repository.NewPlaylistRepo(db), *repository.NewSongsInPlaylistRepo(db)))

	mux.Handle("/create-playlist/", playlistHandler.CreatePlaylist())
	mux.Handle("/delete-playlist/", playlistHandler.DeletePlaylist())
	mux.Handle("/make-public/", playlistHandler.ChangePublicOption(true))
	mux.Handle("/make-private/", playlistHandler.ChangePublicOption(false))
	mux.Handle("/add-song-to-playlist/", playlistHandler.AddSongToPlaylist())
	mux.Handle("/remove-song-from-playlist/", playlistHandler.RemoveSongFromPlaylist())
	mux.Handle("/playlists", playlistHandler.GetPlaylists())
	mux.Handle("/playlists/", playlistHandler.GetPlaylistById())

	services := []string{"http://localhost:9997", "http://localhost:9998", "http://localhost:9999", "http://localhost:9996", "http://localhost:9995"}
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   services,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	log.Println("Playlists service starting on port 9987...")
	if err := http.ListenAndServe(":9987", corsHandler.Handler(mux)); err != nil {
		log.Fatalf("Failed to start Playlists service: %v", err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
