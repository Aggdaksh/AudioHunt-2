package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"Shazam-Vscode/backend/internal/audio/segment"
	"Shazam-Vscode/backend/internal/matching"
)

func main() {
	mode := flag.String("mode", "song-query", "comparison mode: song-query, song-song, query-query")
	flag.Parse()

	if flag.NArg() != 2 {
		log.Fatalf("usage: go run ./cmd/compare-segments --mode <song-query|song-song|query-query> <reference_audio> <query_audio>")
	}

	refPath := flag.Arg(0)
	queryPath := flag.Arg(1)

	refFP, err := fingerprintForMode(*mode, refPath, true)
	if err != nil {
		log.Fatalf("reference fingerprint failed: %v", err)
	}

	queryFP, err := fingerprintForMode(*mode, queryPath, false)
	if err != nil {
		log.Fatalf("query fingerprint failed: %v", err)
	}

	result := matching.SlideHamming(refFP.HashSegments, queryFP.HashSegments)

	fmt.Printf("mode: %s\n", *mode)
	fmt.Printf("reference: %s (%d segments)\n", filepath.Base(refPath), len(refFP.HashSegments))
	fmt.Printf("query: %s (%d segments)\n", filepath.Base(queryPath), len(queryFP.HashSegments))
	fmt.Printf("confidence: %.2f%%\n", result.Confidence*100)
	fmt.Printf("is_match: %v\n", result.IsMatch)
	fmt.Printf("offset: %d\n", result.MatchOffset)
	fmt.Printf("time_in_song: %.2fs\n", result.TimeInSong)
}

func fingerprintForMode(mode, path string, isReference bool) (*segment.Fingerprint, error) {
	switch mode {
	case "song-song":
		return segment.FingerprintSongFile(path)
	case "query-query":
		return segment.FingerprintQueryFile(path)
	case "song-query":
		if isReference {
			return segment.FingerprintSongFile(path)
		}
		return segment.FingerprintQueryFile(path)
	default:
		return nil, fmt.Errorf("unknown mode %q", mode)
	}
}
