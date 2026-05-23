package main

import (
	"fmt"
	"Shazam-Vscode/backend/internal/db"
	"Shazam-Vscode/backend/internal/songs"
)

func main() {
	if err := db.Connect("postgres://shazam_user:pass123@localhost:5434/shazam_clone?sslmode=disable"); err != nil {
		panic(err)
	}
	defer db.Close()

	svc := songs.NewService()

	result, err := svc.Recognize("snippet2.wav")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Match: %s by %s (%.1f%%)\n", result.Song.Title, result.Song.Artist, result.Confidence)
	}
}
