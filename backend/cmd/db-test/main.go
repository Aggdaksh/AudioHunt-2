package main

import (
	"log"

	"Shazam-Vscode/backend/internal/audio/segment"
	"Shazam-Vscode/backend/internal/db"
)

func main() {
	dsn := "postgres://shazam_user:pass123@localhost:5434/shazam_clone?sslmode=disable"
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	filePath := "test.wav"
	fp, err := segment.FingerprintFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	songID, err := db.InsertSong("Test Song", "Test Artist", filePath, fp.Fingerprint, fp.HashSegments)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Inserted song %d with %d segments", songID, len(fp.HashSegments))
}
