package main

import (
	"FrontendService/internal/handler"
	"FrontendService/internal/middleware"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", middleware.VerifyAuthMiddleware(handler.ShowMainPage()))
	mux.Handle("/login/", middleware.VerifyNotAuthMiddleware(handler.LoginHandler()))
	mux.Handle("/register/", middleware.VerifyNotAuthMiddleware(handler.RegisterHandler()))
	mux.Handle("/profile/", middleware.VerifyAuthMiddleware(handler.ShowProfile()))
	mux.Handle("/add-song/", middleware.VerifyAuthMiddleware(handler.AddSong()))
	mux.Handle("/not-an-author/", handler.NotAnAuthor())
	mux.Handle("/become-author/", middleware.VerifyAuthMiddleware(handler.BecomeAuthorPage()))
	mux.Handle("/playlists", middleware.VerifyAuthMiddleware(handler.PlaylistsPage()))
	mux.Handle("/playlists/", middleware.VerifyAuthMiddleware(handler.ShowPlaylist()))

	services := []string{"http://localhost:9997", "http://localhost:9987", "http://localhost:9998", "http://localhost:9999", "http://localhost:9996", "http://localhost:9995"}
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   services,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	log.Println("API Gateway starting on port 9997...")
	if err := http.ListenAndServe(":9997", corsHandler.Handler(mux)); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}
