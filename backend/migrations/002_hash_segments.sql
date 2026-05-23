-- Segment fingerprints (Find_your_tune style) stored per song
ALTER TABLE songs
    ADD COLUMN IF NOT EXISTS fingerprint TEXT,
    ADD COLUMN IF NOT EXISTS hash_segments JSONB;

CREATE INDEX IF NOT EXISTS idx_songs_fingerprint ON songs(fingerprint);
