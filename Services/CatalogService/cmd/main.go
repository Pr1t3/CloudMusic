package main

import (
	"CatalogService/internal/handler"
	"CatalogService/internal/middleware"
	"CatalogService/internal/repository"
	"CatalogService/internal/service"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()
	db, err := openDB("name:password?@/db_name?parseTime=true")
	if err != nil {
		print("Error in opening db")
	}
	defer db.Close()

	catalogService := service.NewCatalogService(*repository.NewSongRepo(db), *repository.NewAuthorRepo(db), *repository.NewGenreRepo(db), *repository.NewSongAuthorRepo(db))

	catalogHandler := handler.NewCatalogHandler(catalogService)

	mux.Handle("/songs", middleware.VerifyAuthMiddleware(catalogHandler.GetSongs()))
	mux.Handle("/songs/", middleware.VerifyAuthMiddleware(catalogHandler.GetSongById()))
	mux.Handle("/genres", middleware.VerifyAuthorMiddleware(middleware.VerifyAuthMiddleware(catalogHandler.GetGenres()), *catalogService))
	mux.Handle("/authors/", middleware.VerifyAuthMiddleware(catalogHandler.GetSongAuthors()))
	mux.Handle("/add-song/", middleware.VerifyAuthorMiddleware(middleware.VerifyAuthMiddleware(catalogHandler.AddSong()), *catalogService))
	mux.Handle("/become-author/", middleware.VerifyNotAuthorMiddleware(middleware.VerifyAuthMiddleware(catalogHandler.BecomeAuthor()), *catalogService))
	mux.Handle("/is-author/", middleware.VerifyAuthMiddleware(catalogHandler.IsAuthor()))
	services := []string{"http://localhost:9997", "http://localhost:9998", "http://localhost:9999", "http://localhost:9996", "http://localhost:9995"}
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   services,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	log.Println("Catalog Service starting on port 9989...")
	if err := http.ListenAndServe(":9989", corsHandler.Handler(mux)); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
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
