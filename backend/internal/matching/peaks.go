package matching

import (
	"math"

	"Shazam-Vscode/backend/internal/audio/dsp"
)

const (
	PeakSampleRate       = 16000
	PeakHopSize          = 512
	PeakMinAlignedHits   = 8
	PeakMinConfidence    = 0.001
	PeakMinLeadOverNext  = 3
	PeakSecondsPerOffset = float64(PeakHopSize) / float64(PeakSampleRate)
	PeakQueryWindowSec   = 8
	PeakQueryStepSec     = 2
)

// MatchPeaks uses a Shazam-style offset histogram over exact hash collisions.
func MatchPeaks(reference, query []dsp.Fingerprint) Result {
	if len(reference) == 0 || len(query) == 0 {
		return Result{IsMatch: false, MatchOffset: -1}
	}

	refByHash := make(map[uint32][]int, len(reference))
	for _, fp := range reference {
		refByHash[fp.Hash] = append(refByHash[fp.Hash], fp.Offset)
	}

	offsetCounts := make(map[int]int)
	bestOffset := -1
	bestCount := 0

	for _, q := range query {
		offsets := refByHash[q.Hash]
		for _, refOffset := range offsets {
			delta := refOffset - q.Offset
			offsetCounts[delta]++
			if offsetCounts[delta] > bestCount {
				bestCount = offsetCounts[delta]
				bestOffset = delta
			}
		}
	}

	secondBest := 0
	for delta, count := range offsetCounts {
		if delta == bestOffset {
			continue
		}
		if count > secondBest {
			secondBest = count
		}
	}

	confidence := float64(bestCount) / float64(len(query))
	isMatch := bestCount >= PeakMinAlignedHits &&
		confidence >= PeakMinConfidence &&
		bestCount >= secondBest+PeakMinLeadOverNext

	return Result{
		IsMatch:     isMatch,
		Confidence:  confidence,
		MatchOffset: bestOffset,
		TimeInSong:  float64(bestOffset) * PeakSecondsPerOffset,
		HitCount:    bestCount,
		RunnerUp:    secondBest,
	}
}

// MatchPeaksWindowed scans overlapping query windows and returns the strongest peak match.
func MatchPeaksWindowed(reference, query []dsp.Fingerprint) Result {
	if len(reference) == 0 || len(query) == 0 {
		return Result{MatchOffset: -1}
	}

	windowFrames := int(math.Round(float64(PeakQueryWindowSec) / PeakSecondsPerOffset))
	stepFrames := int(math.Round(float64(PeakQueryStepSec) / PeakSecondsPerOffset))
	if windowFrames <= 0 {
		windowFrames = 1
	}
	if stepFrames <= 0 {
		stepFrames = 1
	}

	best := MatchPeaks(reference, normalizePeakOffsets(query, 0))
	maxQueryOffset := query[len(query)-1].Offset

	for start := 0; start <= maxQueryOffset; start += stepFrames {
		window := slicePeakWindow(query, start, start+windowFrames)
		if len(window) < PeakMinAlignedHits {
			continue
		}

		result := MatchPeaks(reference, normalizePeakOffsets(window, start))
		if result.HitCount > best.HitCount || (result.HitCount == best.HitCount && result.Confidence > best.Confidence) {
			best = result
		}
	}

	return best
}

func slicePeakWindow(fps []dsp.Fingerprint, start, end int) []dsp.Fingerprint {
	startIdx := 0
	for startIdx < len(fps) && fps[startIdx].Offset < start {
		startIdx++
	}

	endIdx := startIdx
	for endIdx < len(fps) && fps[endIdx].Offset < end {
		endIdx++
	}

	return fps[startIdx:endIdx]
}

func normalizePeakOffsets(fps []dsp.Fingerprint, base int) []dsp.Fingerprint {
	out := make([]dsp.Fingerprint, len(fps))
	for i, fp := range fps {
		out[i] = dsp.Fingerprint{
			Hash:   fp.Hash,
			Offset: fp.Offset - base,
		}
	}
	return out
}
