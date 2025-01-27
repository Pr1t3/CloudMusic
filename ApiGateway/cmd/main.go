package main

import (
	"ApiGateway/internal/handler"
	"ApiGateway/internal/service"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/login/", handler.ProxyHandlerRedirect("http://localhost:9999", "http://localhost:9997"))
	mux.Handle("/register/", handler.ProxyHandlerRedirect("http://localhost:9999", "http://localhost:9997"))
	mux.Handle("/logout/", handler.ProxyHandlerRedirect("http://localhost:9999", "http://localhost:9997/login/"))
	mux.Handle("/change-password", handler.ProxyHandlerRedirect("http://localhost:9999", "http://localhost:9997/profile/"))
	mux.Handle("/upload-photo", handler.ProxyHandler("http://localhost:9999"))
	mux.Handle("/add-song/", handler.ProxyHandlerRedirect("http://localhost:9989", "http://localhost:9997"))
	mux.Handle("/become-author/", handler.ProxyHandlerRedirect("http://localhost:9989", "http://localhost:9997"))
	mux.Handle("/create-playlist/", handler.ProxyHandler("http://localhost:9987"))
	mux.Handle("/add-song-to-playlist/", handler.ProxyHandler("http://localhost:9987"))
	mux.Handle("/delete-playlist/", handler.ProxyHandler("http://localhost:9987"))
	mux.Handle("/remove-song-from-playlist/", handler.ProxyHandler("http://localhost:9987"))
	mux.Handle("/songs/", service.StripSuffix("/", handler.ProxyHandler("http://localhost:9989")))
	mux.Handle("/authors/", handler.ProxyHandler("http://localhost:9989"))
	mux.Handle("/start-song/", handler.ProxyHandler("http://localhost:9988"))

	services := []string{"http://localhost:9997", "http://localhost:9998", "http://localhost:9999", "http://localhost:9996", "http://localhost:9995", "http://localhost:9989", "http://localhost:9988", "http://localhost:9987"}
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   services,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Hash"},
		AllowCredentials: true,
	})

	log.Println("API Gateway starting on port 9998...")
	if err := http.ListenAndServe(":9998", corsHandler.Handler(mux)); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}
