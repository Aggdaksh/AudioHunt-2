-- Songs table
CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    file_path VARCHAR(512) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Fingerprints table
CREATE TABLE IF NOT EXISTS fingerprints (
    id SERIAL PRIMARY KEY,
    hash BIGINT NOT NULL,
    song_id INT NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
    time_offset INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_fingerprints_hash ON fingerprints(hash);
CREATE INDEX IF NOT EXISTS idx_fingerprints_song_id ON fingerprints(song_id);