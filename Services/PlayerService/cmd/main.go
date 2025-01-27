package main

import (
	"PlayerService/internal/handler"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/start-song/", handler.StartSong())

	services := []string{"http://localhost:9997", "http://localhost:9998", "http://localhost:9999", "http://localhost:9996", "http://localhost:9995"}
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   services,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	log.Println("Player Service starting on port 9988...")
	if err := http.ListenAndServe(":9988", corsHandler.Handler(mux)); err != nil {
		log.Fatalf("Failed to start Player Service: %v", err)
	}
}
