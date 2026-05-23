package dsp

import (
	"io"
	"os"

	"github.com/go-audio/wav"
)

// DecodeWav reads a WAV file and returns:
// - samples: slice of float64 audio samples in range [-1.0, 1.0]
// - sampleRate: original sample rate of the file
// - error
//
// The function automatically converts stereo to mono by averaging channels,
// and normalizes 16-bit PCM samples to float64.
func DecodeWav(filePath string) ([]float64, int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)
	if !decoder.IsValidFile() {
		return nil, 0, io.ErrUnexpectedEOF
	}

	// Read entire PCM buffer
	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, 0, err
	}

	numChannels := int(buf.Format.NumChannels)
	sampleRate := buf.Format.SampleRate

	// Convert int PCM to float64 and mix down to mono
	totalSamples := len(buf.Data)
	var samples []float64

	if numChannels == 2 {
		samples = make([]float64, totalSamples/2)
		for i := 0; i < totalSamples; i += 2 {
			// Average left and right channels
			mono := (float64(buf.Data[i]) + float64(buf.Data[i+1])) / 2.0
			samples[i/2] = mono / 32768.0 // Normalize to [-1, 1]
		}
	} else {
		samples = make([]float64, totalSamples)
		for i := 0; i < totalSamples; i++ {
			samples[i] = float64(buf.Data[i]) / 32768.0
		}
	}

	return samples, sampleRate, nil
}
