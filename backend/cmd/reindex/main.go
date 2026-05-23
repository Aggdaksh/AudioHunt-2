package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"Shazam-Vscode/backend/internal/audio/segment"
	"Shazam-Vscode/backend/internal/db"

	"github.com/joho/godotenv"
)

// Rebuilds segment fingerprints for songs missing hash_segments (after algorithm upgrade).
func main() {
	_ = godotenv.Load()
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://shazam_user:pass123@localhost:5434/shazam_clone?sslmode=disable"
	}
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.DB.Query(`SELECT id, title, file_path FROM songs WHERE file_path IS NOT NULL`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title, path string
		if err := rows.Scan(&id, &title, &path); err != nil {
			continue
		}
		fp, err := segment.FingerprintFile(path)
		if err != nil {
			log.Printf("skip %s (%d): %v", title, id, err)
			continue
		}
		segJSON, _ := json.Marshal(fp.HashSegments)
		_, err = db.DB.Exec(
			`UPDATE songs SET fingerprint = $1, hash_segments = $2 WHERE id = $3`,
			fp.Fingerprint, segJSON, id,
		)
		if err != nil {
			log.Printf("update %d failed: %v", id, err)
			continue
		}
		fmt.Printf("reindexed: %s (%d segments)\n", title, len(fp.HashSegments))
	}
}
