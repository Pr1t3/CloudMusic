package repository

import (
	"SearchService/internal/models"
	"database/sql"
)

type TermTrieRepo struct {
	Db *sql.DB
}

func NewTermTrieRepo(db *sql.DB) *TermTrieRepo {
	return &TermTrieRepo{Db: db}
}

func (ttr *TermTrieRepo) InsertTerm(term string) error {
	var parentID *int64
	var err error
	for _, char := range term {
		var nodeID int64
		if parentID == nil {
			err = ttr.Db.QueryRow("SELECT node_id FROM term_trie WHERE letter = ? AND parent_id IS NULL", string(char)).Scan(&nodeID)
		} else {
			err = ttr.Db.QueryRow("SELECT node_id FROM term_trie WHERE letter = ? AND parent_id = ?", string(char), parentID).Scan(&nodeID)
		}
		if err == sql.ErrNoRows {
			res, err := ttr.Db.Exec("INSERT INTO term_trie (letter, parent_id) VALUES (?, ?)", string(char), parentID)
			if err != nil {
				return err
			}
			nodeID, err = res.LastInsertId()
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		parentID = &nodeID
	}

	_, err = ttr.Db.Exec("UPDATE term_trie SET is_end_of_word = TRUE WHERE node_id = ?", *parentID)
	if err != nil {
		return err
	}

	return nil
}

func (ttr *TermTrieRepo) SearchPrefix(prefix string) ([]string, error) {
	var parentID *int
	var err error
	for _, char := range prefix {
		var nodeID int
		if parentID == nil {
			err = ttr.Db.QueryRow("SELECT node_id FROM term_trie WHERE letter = ? AND parent_id IS NULL", string(char)).Scan(&nodeID)
		} else {
			err = ttr.Db.QueryRow("SELECT node_id FROM term_trie WHERE letter = ? AND parent_id = ?", string(char), parentID).Scan(&nodeID)
		}
		if err == sql.ErrNoRows {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		parentID = &nodeID
	}

	return ttr.collectAllWords(*parentID, prefix)
}

func (ttr *TermTrieRepo) collectAllWords(nodeID int, prefix string) ([]string, error) {
	var terms []string
	rows, err := ttr.Db.Query("SELECT node_id, letter, is_end_of_word FROM term_trie WHERE parent_id = ?", nodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var node models.TermTrie
		if err := rows.Scan(&node.NodeId, &node.Letter, &node.IsEndOfWord); err != nil {
			return nil, err
		}
		newPrefix := prefix + node.Letter
		if node.IsEndOfWord {
			terms = append(terms, newPrefix)
		}
		childTerms, err := ttr.collectAllWords(node.NodeId, newPrefix)
		if err != nil {
			return nil, err
		}
		terms = append(terms, childTerms...)
	}

	return terms, nil
}
