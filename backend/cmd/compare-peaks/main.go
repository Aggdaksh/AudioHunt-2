package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"Shazam-Vscode/backend/internal/audio"
	"Shazam-Vscode/backend/internal/audio/dsp"
	"Shazam-Vscode/backend/internal/matching"
)

func main() {
	mode := flag.String("mode", "song-query", "comparison mode: song-query, song-song, query-query")
	flag.Parse()

	if flag.NArg() != 2 {
		log.Fatalf("usage: go run ./cmd/compare-peaks --mode <song-query|song-song|query-query> <reference_audio> <query_audio>")
	}

	refPath := flag.Arg(0)
	queryPath := flag.Arg(1)
	cfg := audio.DefaultConfig()

	refFP, err := fingerprintsForMode(*mode, refPath, cfg, true)
	if err != nil {
		log.Fatalf("reference fingerprint failed: %v", err)
	}
	queryFP, err := fingerprintsForMode(*mode, queryPath, cfg, false)
	if err != nil {
		log.Fatalf("query fingerprint failed: %v", err)
	}

	result := matching.MatchPeaks(refFP, queryFP)
	if *mode != "song-song" {
		result = matching.MatchPeaksWindowed(refFP, queryFP)
	}

	fmt.Printf("mode: %s\n", *mode)
	fmt.Printf("reference: %s (%d fingerprints)\n", filepath.Base(refPath), len(refFP))
	fmt.Printf("query: %s (%d fingerprints)\n", filepath.Base(queryPath), len(queryFP))
	fmt.Printf("aligned_hits: %d\n", result.HitCount)
	fmt.Printf("runner_up_hits: %d\n", result.RunnerUp)
	fmt.Printf("confidence: %.4f\n", result.Confidence)
	fmt.Printf("is_match: %v\n", result.IsMatch)
	fmt.Printf("offset: %d\n", result.MatchOffset)
	fmt.Printf("time_in_song: %.2fs\n", result.TimeInSong)
}

func fingerprintsForMode(mode, path string, cfg audio.FingerprintConfig, isReference bool) ([]dsp.Fingerprint, error) {
	switch mode {
	case "song-song":
		return audio.ExtractFingerprints(path, cfg)
	case "query-query":
		return audio.ExtractQueryFingerprints(path, cfg)
	case "song-query":
		if isReference {
			return audio.ExtractFingerprints(path, cfg)
		}
		return audio.ExtractQueryFingerprints(path, cfg)
	default:
		return nil, fmt.Errorf("unknown mode %q", mode)
	}
}
