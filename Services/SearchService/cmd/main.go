package main

import (
	"SearchService/internal/handler"
	"SearchService/internal/middleware"
	"SearchService/internal/repository"
	"SearchService/internal/service"
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

	searchHandler := handler.NewSearchHandler(*service.NewSearchService(*repository.NewSearchRepo(db), *repository.NewTermTrieRepo(db)))

	mux.Handle("/add-term/", middleware.VerifyAuthMiddleware(searchHandler.InsertTerm()))
	mux.Handle("/search/", middleware.VerifyAuthMiddleware(searchHandler.SearchByPrefix()))
	mux.Handle("/get-entities/", middleware.VerifyAuthMiddleware(searchHandler.GetEntityByTerm()))

	services := []string{"http://localhost:9997", "http://localhost:9998", "http://localhost:9999", "http://localhost:9996", "http://localhost:9995"}
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   services,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	log.Println("Search Service starting on port 9986...")
	if err := http.ListenAndServe(":9986", corsHandler.Handler(mux)); err != nil {
		log.Fatalf("Failed to start Search Service: %v", err)
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
