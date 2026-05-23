package dsp

import (
	"encoding/json"
	"os"
	"time"
)

func logFingerprintDebug(hypothesisID, location, message string, data map[string]interface{}) {
	logPath := os.Getenv("AUDIOHUNT_DEBUG_LOG")
	if logPath == "" {
		return
	}
	payload := map[string]interface{}{
		"sessionId":    os.Getenv("AUDIOHUNT_DEBUG_SESSION"),
		"runId":        os.Getenv("AUDIOHUNT_DEBUG_RUN"),
		"hypothesisId": hypothesisID,
		"location":     location,
		"message":      message,
		"data":         data,
		"timestamp":    time.Now().UnixMilli(),
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.Write(append(b, '\n'))
}

// Fingerprint is a hash + the absolute time offset (in frames) of the anchor peak.
type Fingerprint struct {
	Hash   uint32
	Offset int
}

const (
	freqBinQuantization   = 4
	timeDeltaQuantization = 2
)

// GenerateFingerprints creates combinatorial hashes from a list of peaks.
// Uses the classic Shazam technique:
//   - Sort peaks by time.
//   - For each peak (anchor), pair it with up to fanValue next peaks (targets).
//   - For each pair, compute a 32-bit hash from anchor frequency, target frequency,
//     and time delta.
//
// Parameters:
//   - peaks: slice of Peak structs
//   - fanValue: how many target peaks to pair with each anchor (e.g., 15)
//   - maxDelta: maximum allowed time difference in frames (e.g., 200)
//
// Returns a slice of Fingerprints.
func GenerateFingerprints(peaks []Peak, fanValue, maxDelta int) []Fingerprint {
	if len(peaks) < 2 {
		return nil
	}
	skippedForDelta := 0

	// Sort peaks by time (ascending)
	sorted := make([]Peak, len(peaks))
	copy(sorted, peaks)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Time < sorted[i].Time {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	fingerprints := []Fingerprint{}

	for i := 0; i < len(sorted); i++ {
		anchor := sorted[i]
		// Pair with next fanValue peaks
		for j := 1; j <= fanValue && i+j < len(sorted); j++ {
			target := sorted[i+j]
			delta := target.Time - anchor.Time
			if delta <= 0 || delta > maxDelta {
				continue
			}
			if delta > 255 {
				skippedForDelta++
				continue
			}

			anchorBin := anchor.FreqBin / freqBinQuantization
			targetBin := target.FreqBin / freqBinQuantization
			deltaBucket := delta / timeDeltaQuantization

			// 32-bit hash:
			// Quantized bins make microphone captures less brittle to small spectral drift.
			hash := (uint32(anchorBin) << 16) |
				(uint32(targetBin) << 8) |
				uint32(deltaBucket)

			fingerprints = append(fingerprints, Fingerprint{
				Hash:   hash,
				Offset: anchor.Time,
			})
		}
	}
	// #region agent log
	logFingerprintDebug("H3", "fingerprint.go:GenerateFingerprints", "Fingerprint generation stats", map[string]interface{}{"peakCount": len(peaks), "fingerprintCount": len(fingerprints), "skippedDeltaOutOf8Bit": skippedForDelta, "maxDelta": maxDelta, "freqQuant": freqBinQuantization, "deltaQuant": timeDeltaQuantization})
	// #endregion
	return fingerprints
}
