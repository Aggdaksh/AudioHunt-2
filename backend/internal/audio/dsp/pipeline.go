package dsp

type Config struct {
	FanValue   int
	WindowSize int
	HopSize    int
}

func DefaultConfig() Config {
	return Config{
		FanValue:   5,
		WindowSize: 2048,
		HopSize:    512,
	}
}

func ExtractFingerprints(filePath string, cfg Config) ([]Fingerprint, error) {

	// Step 1: WAV
	samples, sampleRate, err := DecodeWav(filePath)
	if err != nil {
		return nil, err
	}

	// Step 2: Spectrogram
	spec := GenerateSpectrogram(samples, sampleRate, cfg.WindowSize, cfg.HopSize)

	// Step 3: Peaks
	peaks := FindPeaks(spec, 10, 20.0)

	// Step 4: Fingerprints ✅ FIXED
	fingerprints := GenerateFingerprints(peaks, cfg.FanValue, 200)

	return fingerprints, nil

}
