use CloudMusic_CatalogService;

drop table authors;
drop table genres;
drop table albums;
drop table songs;
drop table song_authors;

CREATE TABLE authors (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE genres (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE songs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    duration int NOT NULL COMMENT "in seconds",
    size BIGINT not null,
    genre_id INT,
    file_path VARCHAR(255) NOT NULL,
    added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE SET NULL
);

CREATE TABLE song_authors (
    song_id INT NOT NULL,
    author_id INT NOT NULL,
    PRIMARY KEY (song_id, author_id),
    FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE CASCADE
);
