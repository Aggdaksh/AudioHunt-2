//go:build ignore

package main

import (
	"Shazam-Vscode/backend/internal/db"
	"fmt"
)

func main() {
	if err := db.Connect("postgres://shazam_user:pass123@localhost:5434/shazam_clone?sslmode=disable"); err != nil {
		panic(err)
	}
	defer db.Close()
	songs, err := db.GetSongsWithSegments()
	if err != nil {
		panic(err)
	}
	for _, s := range songs {
		fmt.Printf("ID: %d, Title: %s, Segments: %d\n", s.ID, s.Title, len(s.HashSegments))
	}
}
