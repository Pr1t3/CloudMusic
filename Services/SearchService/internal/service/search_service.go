package service

import (
	"SearchService/internal/models"
	"SearchService/internal/repository"
)

type SearchService struct {
	IndexRepo    *repository.SearchIndexRepo
	TermTrieRepo *repository.TermTrieRepo
}

func NewSearchService(ii repository.SearchIndexRepo, tt repository.TermTrieRepo) *SearchService {
	return &SearchService{IndexRepo: &ii, TermTrieRepo: &tt}
}

func (s *SearchService) InsertTerm(term string, entityId int, entityType string) error {
	err := s.IndexRepo.InsertTerm(term, entityId, entityType)
	if err != nil {
		return err
	}
	return s.TermTrieRepo.InsertTerm(term)
}

func (s *SearchService) SearchPrefix(prefix string) ([]string, error) {
	return s.TermTrieRepo.SearchPrefix(prefix)
}

func (s *SearchService) GetEntities(term string) ([]models.SearchIndex, error) {
	return s.IndexRepo.GetEntities(term)
}
