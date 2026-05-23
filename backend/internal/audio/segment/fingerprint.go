package segment

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/cmplx"
	"sort"
	"strings"

	"github.com/madelynnblue/go-dsp/fft"
)

// Fingerprint holds per-window spectral hashes (Find_your_tune style).
type Fingerprint struct {
	Fingerprint  string
	HashSegments []string
}

const (
	SampleRate = 22050
	WindowSize = 2048
	HopSize    = 512
)

func Generate(samples []float64) (*Fingerprint, error) {
	if len(samples) < WindowSize {
		return nil, fmt.Errorf("insufficient audio samples")
	}

	var hashSegments []string
	for i := 0; i < len(samples)-WindowSize; i += HopSize {
		window := samples[i : i+WindowSize]
		if rms(window) < 0.00001 {
			continue
		}

		windowed := make([]float64, WindowSize)
		for j, sample := range window {
			w := 0.54 - 0.46*math.Cos(2*math.Pi*float64(j)/float64(WindowSize-1))
			windowed[j] = sample * w
		}

		spectrum := fft.FFTReal(windowed)
		hash := createRobustHash(spectrum)
		if hash != "" {
			hashSegments = append(hashSegments, hash)
		}
	}

	if len(hashSegments) < 1 {
		return nil, fmt.Errorf("too few hash segments generated: %d", len(hashSegments))
	}

	combinedHash := strings.Join(hashSegments, "")
	finalHash := sha256.Sum256([]byte(combinedHash))

	return &Fingerprint{
		Fingerprint:  hex.EncodeToString(finalHash[:]),
		HashSegments: hashSegments,
	}, nil
}

func createRobustHash(spectrum []complex128) string {
	n := len(spectrum)
	if n > 1024 {
		n = 1024 // Match original 1024 positive frequency bins
	}
	
	mags := make([]float64, n)
	const eps = 1e-12
	for i := 0; i < n; i++ {
		mags[i] = cmplx.Abs(spectrum[i]) + eps
	}

	const numBands = 16
	bandSize := n / numBands
	bands := make([]float64, numBands)

	for b := 0; b < numBands; b++ {
		start := b * bandSize
		end := start + bandSize
		if b == numBands-1 || end > n {
			end = n
		}
		var sum float64
		for i := start; i < end; i++ {
			sum += mags[i]
		}
		bands[b] = 20.0 * math.Log10(sum/float64(end-start))
	}

	m := median(bands)
	for b := range bands {
		bands[b] -= m
	}

	const lo, hi = -24.0, 24.0
	hash := make([]byte, numBands)
	hexDigits := []byte("0123456789abcdef")
	for b := 0; b < numBands; b++ {
		v := bands[b]
		if v < lo {
			v = lo
		}
		if v > hi {
			v = hi
		}
		q := int(math.Round((v - lo) / (hi - lo) * 15.0))
		hash[b] = hexDigits[q]
	}
	return string(hash)
}

func median(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	tmp := append([]float64(nil), xs...)
	sort.Float64s(tmp)
	n := len(tmp)
	if n%2 == 1 {
		return tmp[n/2]
	}
	return 0.5 * (tmp[n/2-1] + tmp[n/2])
}



func rms(samples []float64) float64 {
	var sum float64
	for _, sample := range samples {
		sum += sample * sample
	}
	return math.Sqrt(sum / float64(len(samples)))
}
