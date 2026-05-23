package dsp

import (
	"math"
	"math/cmplx"

	"github.com/madelynnblue/go-dsp/fft"
)

// GenerateSpectrogram creates a magnitude spectrogram (in dB) from audio samples.
// Parameters:
//   - samples: input audio samples
//   - sampleRate: original sample rate (unused here, kept for future resampling)
//   - windowSize: size of FFT window (e.g., 2048)
//   - hopSize: number of samples between successive windows (e.g., 512)
//
// Returns a 2D slice: spectrogram[frequencyBin][timeFrame] = dB magnitude
func GenerateSpectrogram(samples []float64, sampleRate int, windowSize, hopSize int) [][]float64 {
	if len(samples) < windowSize {
		return nil
	}

	numFrames := (len(samples)-windowSize)/hopSize + 1
	freqBins := windowSize/2 + 1

	// Initialize spectrogram: rows = frequency bins, columns = time frames
	spec := make([][]float64, freqBins)
	for i := range spec {
		spec[i] = make([]float64, numFrames)
	}

	// Precompute Hann window
	hann := make([]float64, windowSize)
	for i := 0; i < windowSize; i++ {
		hann[i] = 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(windowSize-1)))
	}

	window := make([]float64, windowSize)

	for frame := 0; frame < numFrames; frame++ {
		start := frame * hopSize

		// Apply window
		for i := 0; i < windowSize; i++ {
			window[i] = samples[start+i] * hann[i]
		}

		// Compute FFT
		complexFFT := fft.FFTReal(window)

		// Convert to dB magnitude
		for bin := 0; bin < freqBins; bin++ {
			mag := cmplx.Abs(complexFFT[bin])
			if mag < 1e-10 {
				mag = 1e-10
			}
			spec[bin][frame] = 20 * math.Log10(mag)
		}
	}

	return spec
}
