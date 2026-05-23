//go:build ignore

package main

import (
	"Shazam-Vscode/backend/internal/audio/segment"
	"fmt"
	"math"
	"math/rand"
)

func HammingHex(a, b string) int {
	d := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			d++
		}
	}
	return d
}

func main() {
	N := 22050 * 20 // 20 seconds
	orig := make([]float64, N)
	mic := make([]float64, N)

	for i := 0; i < N; i++ {
		orig[i] = math.Sin(2*math.Pi*440*float64(i)/22050) +
			math.Sin(2*math.Pi*880*float64(i)/22050) +
			math.Sin(2*math.Pi*2000*float64(i)/22050)

		mic[i] = 0.1*math.Sin(2*math.Pi*440*float64(i)/22050) +
			math.Sin(2*math.Pi*880*float64(i)/22050) +
			2.0*math.Sin(2*math.Pi*2000*float64(i)/22050) +
			(rand.Float64()-0.5)*0.5
	}

	fpOrig, _ := segment.Generate(orig)
	fpMic, _ := segment.Generate(mic)

	ref := fpOrig.HashSegments
	qry := fpMic.HashSegments

	maxNibbles := []int{8, 9, 10, 11, 12}

	for _, maxNibbleMismatches := range maxNibbles {
		bestConf := 0.0
		for offset := 0; offset <= len(ref)-len(qry); offset++ {
			matches := 0
			for i := 0; i < len(qry); i++ {
				// NO +/- 1 SHIFT!
				if HammingHex(ref[offset+i], qry[i]) <= maxNibbleMismatches {
					matches++
				}
			}
			conf := float64(matches) / float64(len(qry))
			if conf > bestConf {
				bestConf = conf
			}
		}
		fmt.Printf("Max Mismatches = %d, Confidence = %.2f%%\n", maxNibbleMismatches, bestConf*100)
	}
}
