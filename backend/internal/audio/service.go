package audio

import (
	"fmt"
	"os"
	"os/exec"

	"Shazam-Vscode/backend/internal/audio/dsp"
)

// FingerprintConfig holds parameters for the fingerprinting pipeline.
type FingerprintConfig struct {
	// Spectrogram parameters
	WindowSize int
	HopSize    int

	// Peak picking
	PeakNeighborhood int
	PeakMinAmplitude float64

	// Hashing
	FanValue int
	MaxDelta int
}

// DefaultConfig returns sensible values used by Shazam-like systems.
func DefaultConfig() FingerprintConfig {
	return FingerprintConfig{
		WindowSize:       2048,
		HopSize:          512,
		PeakNeighborhood: 10,
		PeakMinAmplitude: 20.0,
		FanValue:         15,
		MaxDelta:         200,
	}
}

func convertToWav(filePath, outPath, audioFilter string) error {
	args := []string{"-y", "-i", filePath}
	if audioFilter != "" {
		args = append(args, "-af", audioFilter)
	}
	args = append(args, "-ac", "1", "-ar", "16000", outPath)
	out, err := exec.Command("ffmpeg", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg conversion failed: %w, output: %s", err, string(out))
	}
	return nil
}

func ExtractFingerprints(filePath string, cfg FingerprintConfig) ([]dsp.Fingerprint, error) {
	return extractFingerprintsFromFile(filePath, cfg, "")
}

// ExtractQueryFingerprints preprocesses mic/noisy snippets before fingerprinting.
// Enrollment keeps ExtractFingerprints (clean source files).
func ExtractQueryFingerprints(filePath string, cfg FingerprintConfig) ([]dsp.Fingerprint, error) {
	filter := "highpass=f=80,lowpass=8000,afftdn=nf=-12,loudnorm=I=-16:TP=-1.5:LRA=11"
	fps, err := extractFingerprintsFromFile(filePath, cfg, filter)
	if err != nil {
		filter = "highpass=f=80,lowpass=8000,loudnorm=I=-16:TP=-1.5:LRA=11"
		return extractFingerprintsFromFile(filePath, cfg, filter)
	}
	return fps, nil
}

func extractFingerprintsFromFile(filePath string, cfg FingerprintConfig, audioFilter string) ([]dsp.Fingerprint, error) {
	tmpWav, err := os.CreateTemp("", "converted-*.wav")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpWavPath := tmpWav.Name()
	tmpWav.Close()
	defer os.Remove(tmpWavPath)

	if err := convertToWav(filePath, tmpWavPath, audioFilter); err != nil {
		return nil, err
	}

	// Decode
	samples, sr, err := dsp.DecodeWav(tmpWavPath)
	if err != nil {
		return nil, err
	}

	// Spectrogram
	spec := dsp.GenerateSpectrogram(samples, sr, cfg.WindowSize, cfg.HopSize)
	if spec == nil {
		return nil, nil
	}

	// Peaks
	peaks := dsp.FindPeaks(spec, cfg.PeakNeighborhood, cfg.PeakMinAmplitude)

	// Fingerprints
	fingerprints := dsp.GenerateFingerprints(peaks, cfg.FanValue, cfg.MaxDelta)

	return fingerprints, nil
}
