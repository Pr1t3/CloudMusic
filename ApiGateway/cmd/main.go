package main

import (
	"ApiGateway/internal/handler"
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

	services := []string{"http://localhost:9997", "http://localhost:9998", "http://localhost:9999", "http://localhost:9996", "http://localhost:9995"}
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   services,                                 // Разрешаем только домен фронтенда
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"}, // Разрешаем методы
		AllowedHeaders:   []string{"Content-Type", "Hash"},         // Разрешаем заголовок Content-Type
		AllowCredentials: true,                                     // Разрешаем куки
	})

	log.Println("API Gateway starting on port 9998...")
	if err := http.ListenAndServe(":9998", corsHandler.Handler(mux)); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}
