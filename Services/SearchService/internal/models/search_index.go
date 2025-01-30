package models

type SearchIndex struct {
	Term       string `json:"term"`
	EntityId   int    `json:"entity_id"`
	EntityType string `json:"entity_type"`
}
