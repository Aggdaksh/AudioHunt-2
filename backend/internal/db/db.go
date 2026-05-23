package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect establishes a connection to PostgreSQL.
func Connect(dsn string) error {
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	log.Println("Connected to database")
	return nil
}

// Close closes the database connection.
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// Song represents a row in the songs table.
type Song struct {
	ID           int
	Title        string
	Artist       string
	FilePath     string
	Fingerprint  string
	HashSegments []string
}

type PeakFingerprint struct {
	Hash   uint32
	Offset int
}

// InsertSong adds a song with segment fingerprint and returns its ID.
func InsertSong(title, artist, filePath, fingerprint string, hashSegments []string) (int, error) {
	segmentsJSON, err := json.Marshal(hashSegments)
	if err != nil {
		return 0, err
	}

	var id int
	err = DB.QueryRow(
		`INSERT INTO songs (title, artist, file_path, fingerprint, hash_segments)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		title, artist, filePath, fingerprint, segmentsJSON,
	).Scan(&id)

	return id, err
}

// InsertSongWithPeakFingerprints adds a song and its peak fingerprints in one transaction.
func InsertSongWithPeakFingerprints(
	title, artist, filePath, fingerprint string,
	hashSegments []string,
	peakFingerprints []PeakFingerprint,
) (int, error) {
	segmentsJSON, err := json.Marshal(hashSegments)
	if err != nil {
		return 0, err
	}

	tx, err := DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var id int
	if err := tx.QueryRow(
		`INSERT INTO songs (title, artist, file_path, fingerprint, hash_segments)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		title, artist, filePath, fingerprint, segmentsJSON,
	).Scan(&id); err != nil {
		return 0, err
	}

	if len(peakFingerprints) > 0 {
		stmt, err := tx.Prepare(`INSERT INTO fingerprints (hash, song_id, time_offset) VALUES ($1, $2, $3)`)
		if err != nil {
			return 0, err
		}
		defer stmt.Close()

		for _, fp := range peakFingerprints {
			if _, err := stmt.Exec(int64(fp.Hash), id, fp.Offset); err != nil {
				return 0, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func InsertPeakFingerprints(songID int, peakFingerprints []PeakFingerprint) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM fingerprints WHERE song_id = $1`, songID); err != nil {
		return err
	}

	if len(peakFingerprints) > 0 {
		stmt, err := tx.Prepare(`INSERT INTO fingerprints (hash, song_id, time_offset) VALUES ($1, $2, $3)`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, fp := range peakFingerprints {
			if _, err := stmt.Exec(int64(fp.Hash), songID, fp.Offset); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func GetPeakFingerprints(songID int) ([]PeakFingerprint, error) {
	rows, err := DB.Query(`
		SELECT hash, time_offset
		FROM fingerprints
		WHERE song_id = $1
		ORDER BY time_offset
	`, songID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fingerprints []PeakFingerprint
	for rows.Next() {
		var hash int64
		var offset int
		if err := rows.Scan(&hash, &offset); err != nil {
			return nil, err
		}
		fingerprints = append(fingerprints, PeakFingerprint{
			Hash:   uint32(hash),
			Offset: offset,
		})
	}

	return fingerprints, rows.Err()
}

// GetSongsWithSegments returns enrolled songs that have segment fingerprints.
func GetSongsWithSegments() ([]Song, error) {
	rows, err := DB.Query(`
		SELECT id, title, artist, file_path, fingerprint, hash_segments
		FROM songs
		WHERE hash_segments IS NOT NULL AND jsonb_array_length(hash_segments) > 0
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []Song
	for rows.Next() {
		var s Song
		var segmentsJSON []byte
		if err := rows.Scan(&s.ID, &s.Title, &s.Artist, &s.FilePath, &s.Fingerprint, &segmentsJSON); err != nil {
			continue
		}
		_ = json.Unmarshal(segmentsJSON, &s.HashSegments)
		if len(s.HashSegments) > 0 {
			songs = append(songs, s)
		}
	}
	return songs, nil
}
