package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"Shazam-Vscode/backend/internal/audio"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path_to_wav_file>\n", os.Args[0])
		os.Exit(1)
	}

	filePath := os.Args[1]
	cfg := audio.DefaultConfig()

	fingerprints, err := audio.ExtractFingerprints(filePath, cfg)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Print as JSON for readability
	output, _ := json.MarshalIndent(fingerprints, "", "  ")
	fmt.Println(string(output))

	log.Printf("Generated %d fingerprints from %s\n", len(fingerprints), filePath)
}
