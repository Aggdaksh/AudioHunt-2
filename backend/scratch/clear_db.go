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

	_, err := db.DB.Exec("DELETE FROM songs")
	if err != nil {
		panic(err)
	}

	// Also reset the auto-increment counter
	_, err = db.DB.Exec("ALTER SEQUENCE songs_id_seq RESTART WITH 1")
	if err != nil {
		fmt.Println("Warning: Could not reset sequence:", err)
	}

	fmt.Println("Database successfully cleared!")
}
