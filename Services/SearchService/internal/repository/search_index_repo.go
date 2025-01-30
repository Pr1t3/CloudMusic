package repository

import (
	"SearchService/internal/models"
	"database/sql"
)

type SearchIndexRepo struct {
	Db *sql.DB
}

func NewSearchRepo(db *sql.DB) *SearchIndexRepo {
	return &SearchIndexRepo{Db: db}
}

func (r *SearchIndexRepo) GetEntities(term string) ([]models.SearchIndex, error) {
	query := `SELECT * FROM search_index WHERE term = ?`
	rows, err := r.Db.Query(query, term)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var searchIndexes []models.SearchIndex
	for rows.Next() {
		searchIndex := models.SearchIndex{}
		if err := rows.Scan(&searchIndex.Term, &searchIndex.EntityId, &searchIndex.EntityType); err != nil {
			return nil, err
		}
		searchIndexes = append(searchIndexes, searchIndex)
	}
	return searchIndexes, nil
}

func (r *SearchIndexRepo) InsertTerm(term string, entityId int, entityType string) error {
	_, err := r.Db.Exec("INSERT INTO search_index (term, entity_id, entity_type) VALUES (?, ?, ?)", term, entityId, entityType)
	return err
}
