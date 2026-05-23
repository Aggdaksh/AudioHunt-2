package segment

import (
	"fmt"
	"math"
	"os"
	"os/exec"
)

var songFilters = []string{
	"",
	"highpass=f=80,lowpass=f=8000",
}

var queryFilters = []string{
	"highpass=f=80,lowpass=f=6000,afftdn=nf=-18,loudnorm=I=-16:TP=-1.5:LRA=11",
	"highpass=f=80,lowpass=f=6000,loudnorm=I=-16:TP=-1.5:LRA=11",
	"highpass=f=80,lowpass=f=8000",
	"",
}

func convertWithFfmpeg(filePath, outPath, filter string) error {
	args := []string{"-y", "-i", filePath}
	if filter != "" {
		args = append(args, "-af", filter)
	}
	args = append(args, "-ac", "1", "-ar", fmt.Sprintf("%d", SampleRate), outPath)
	out, err := exec.Command("ffmpeg", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w, output: %s", err, string(out))
	}
	return nil
}

// DecodeFile converts any audio file to mono float64 samples at 22050 Hz.
func decodeFileWithFilters(filePath string, filters []string) ([]float64, error) {
	tmpWav, err := os.CreateTemp("", "segment-*.wav")
	if err != nil {
		return nil, err
	}
	tmpPath := tmpWav.Name()
	tmpWav.Close()
	defer os.Remove(tmpPath)

	var convErr error
	for _, f := range filters {
		if err := convertWithFfmpeg(filePath, tmpPath, f); err == nil {
			convErr = nil
			break
		} else {
			convErr = err
		}
	}
	if convErr != nil {
		return nil, convErr
	}

	cmd := exec.Command("ffmpeg", "-i", tmpPath, "-f", "f64le", "-acodec", "pcm_f64le", "-ac", "1", "-ar", fmt.Sprintf("%d", SampleRate), "-")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to decode audio: %w", err)
	}

	samples := make([]float64, len(output)/8)
	for i := 0; i < len(samples); i++ {
		bits := uint64(0)
		for j := 0; j < 8; j++ {
			bits |= uint64(output[i*8+j]) << (j * 8)
		}
		samples[i] = math.Float64frombits(bits)
	}
	return samples, nil
}

func DecodeSongFile(filePath string) ([]float64, error) {
	return decodeFileWithFilters(filePath, songFilters)
}

func DecodeQueryFile(filePath string) ([]float64, error) {
	return decodeFileWithFilters(filePath, queryFilters)
}

// DecodeFile is kept for backward compatibility with older tools and scripts.
func DecodeFile(filePath string) ([]float64, error) {
	return DecodeSongFile(filePath)
}

func fingerprintDecodedSamples(samples []float64) (*Fingerprint, error) {
	return Generate(samples)
}

func FingerprintSongFile(filePath string) (*Fingerprint, error) {
	samples, err := DecodeSongFile(filePath)
	if err != nil {
		return nil, err
	}
	return fingerprintDecodedSamples(samples)
}

func FingerprintQueryFile(filePath string) (*Fingerprint, error) {
	samples, err := DecodeQueryFile(filePath)
	if err != nil {
		return nil, err
	}
	return fingerprintDecodedSamples(samples)
}

// FingerprintFile is kept for backward compatibility with older tools and scripts.
func FingerprintFile(filePath string) (*Fingerprint, error) {
	return FingerprintSongFile(filePath)
}
