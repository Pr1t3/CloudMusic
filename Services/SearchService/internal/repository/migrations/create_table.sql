use CloudMusic_SearchService;

DROP TABLE search_index;
DROP TABLE term_trie;

CREATE TABLE search_index (
    term VARCHAR(255) NOT NULL,
    entity_id INT NOT NULL,
    entity_type VARCHAR(255) NOT NULL,
    PRIMARY KEY (term, entity_id, entity_type)
);

CREATE TABLE term_trie (
    node_id INT PRIMARY KEY AUTO_INCREMENT,
    letter CHAR(1) NOT NULL,
    parent_id INT,
    is_end_of_word TINYINT(1) DEFAULT 0,
    FOREIGN KEY (parent_id) REFERENCES term_trie(node_id)
);
