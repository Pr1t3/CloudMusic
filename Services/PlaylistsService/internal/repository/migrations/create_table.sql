use CloudMusic_Playlists;

drop table playlists;
drop table songs_in_playlists;

CREATE TABLE songs_in_playlists (
    playlist_id INT NOT NULL,
    song_id INT NOT NULL,
    song_order INT DEFAULT 0,
    PRIMARY KEY (playlist_id, song_id),
    FOREIGN KEY (playlist_id) REFERENCES playlists(id) ON DELETE CASCADE
);

CREATE TABLE playlists (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    user_id INT NOT NULL,
    is_public BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
