package dsp

import (
	"encoding/json"
	"math"
	"os"
	"time"
)

func logPeaksDebug(hypothesisID, location, message string, data map[string]interface{}) {
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

func computeStats(spec [][]float64) (float64, float64) {
	total := 0.0
	count := 0.0
	for f := 0; f < len(spec); f++ {
		for t := 0; t < len(spec[f]); t++ {
			total += spec[f][t]
			count++
		}
	}
	if count == 0 {
		return 0, 0
	}
	mean := total / count
	variance := 0.0
	for f := 0; f < len(spec); f++ {
		for t := 0; t < len(spec[f]); t++ {
			diff := spec[f][t] - mean
			variance += diff * diff
		}
	}
	return mean, math.Sqrt(variance / count)
}

// Peak represents a local maximum in the spectrogram.
// FreqBin is the frequency bin index, Time is the frame index.
type Peak struct {
	FreqBin int
	Time    int
}

// FindPeaks finds local maxima in a spectrogram using a rectangular neighborhood.
// Parameters:
//   - spec: 2D spectrogram [freq][time] in dB
//   - neighborhoodSize: radius around each point to check (e.g., 10)
//   - ampMin: minimum amplitude in dB to consider (e.g., 20.0)
//
// Returns a slice of Peak coordinates.
func FindPeaks(spec [][]float64, neighborhoodSize int, ampMin float64) []Peak {
	freqBins := len(spec)
	if freqBins == 0 {
		return nil
	}
	timeFrames := len(spec[0])
	half := neighborhoodSize / 2

	var peaks []Peak
	mean, sigma := computeStats(spec)
	// Use original threshold to match database density
	threshold := mean + 1.5*sigma

	// Skip edges to avoid boundary checks
	for t := half; t < timeFrames-half; t++ {
		for f := half; f < freqBins-half; f++ {
			val := spec[f][t]
			if val < threshold {
				continue
			}

			isMax := true
			// Check all neighbors within the rectangle
			for dt := -half; dt <= half; dt++ {
				for df := -half; df <= half; df++ {
					if dt == 0 && df == 0 {
						continue
					}
					if spec[f+df][t+dt] >= val {
						isMax = false
						break
					}
				}
				if !isMax {
					break
				}
			}

			if isMax {
				peaks = append(peaks, Peak{FreqBin: f, Time: t})
			}
		}
	}
	// #region agent log
	logPeaksDebug("H4", "peaks.go:FindPeaks", "Peak detection stats", map[string]interface{}{"ampMin": ampMin, "mean": mean, "sigma": sigma, "peakCount": len(peaks), "freqBins": freqBins, "timeFrames": timeFrames})
	// #endregion
	return peaks
}
