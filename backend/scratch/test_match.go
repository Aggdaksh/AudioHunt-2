//go:build ignore

package main

import (
	"Shazam-Vscode/backend/internal/db"
	"Shazam-Vscode/backend/internal/matching"
	"fmt"
)

func main() {
	if err := db.Connect("postgres://shazam_user:pass123@localhost:5434/shazam_clone?sslmode=disable"); err != nil {
		panic(err)
	}
	defer db.Close()
	songs, _ := db.GetSongsWithSegments()

	// Find Shape of You (ID 6) and (ID 7)
	var s6, s7 db.Song
	for _, s := range songs {
		if s.ID == 6 {
			s6 = s
		}
		if s.ID == 7 {
			s7 = s
		}
	}

	if len(s6.HashSegments) == 0 || len(s7.HashSegments) == 0 {
		fmt.Println("Missing segments")
		return
	}

	// Truncate s7 to simulate a 10 second snippet (430 segments)
	snippet := s7.HashSegments[1000:1430]

	res := matching.SlideHamming(s6.HashSegments, snippet)
	fmt.Printf("Self Match Confidence: %.2f%%\n", res.Confidence*100)
}
