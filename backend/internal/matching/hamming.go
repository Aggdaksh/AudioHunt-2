package matching

import (
	"math"

	"Shazam-Vscode/backend/internal/audio/segment"
)

const (
	MatchThreshold      = 0.25
	MaxNibbleMismatches = 8

	TimePerSegmentSeconds = float64(segment.HopSize) / float64(segment.SampleRate)
)

type Result struct {
	IsMatch     bool
	Confidence  float64 // 0.0–1.0, fraction of query segments matched
	MatchOffset int
	TimeInSong  float64
	HitCount    int
	RunnerUp    int
}

func HammingHex(a, b string) int {
	if len(a) != len(b) {
		return math.MaxInt32
	}
	d := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			d++
		}
	}
	return d
}

// SlideHamming slides the query fingerprint across a reference (Find_your_tune logic).
func SlideHamming(reference, query []string) Result {
	if len(query) == 0 || len(reference) == 0 {
		return Result{IsMatch: false, Confidence: 0, MatchOffset: -1}
	}

	best := Result{IsMatch: false, Confidence: 0, MatchOffset: -1}
	maxOffset := len(reference) - len(query)
	if maxOffset < 0 {
		maxOffset = 0
	}

	for offset := 0; offset <= maxOffset; offset++ {
		matches := 0
		for i := 0; i < len(query) && offset+i < len(reference); i++ {
			if HammingHex(reference[offset+i], query[i]) <= MaxNibbleMismatches {
				matches++
			}
		}
		conf := float64(matches) / float64(len(query))
		if conf > best.Confidence {
			best = Result{
				IsMatch:     conf >= MatchThreshold,
				Confidence:  conf,
				MatchOffset: offset,
				TimeInSong:  float64(offset) * TimePerSegmentSeconds,
				HitCount:    matches,
			}
		}
	}

	return best
}
