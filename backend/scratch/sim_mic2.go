//go:build ignore

package main

import (
	"Shazam-Vscode/backend/internal/audio/segment"
	"Shazam-Vscode/backend/internal/matching"
	"fmt"
	"math"
	"math/rand"
)

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

	res := matching.SlideHamming(fpOrig.HashSegments, fpMic.HashSegments)
	fmt.Printf("Simulated Mic Confidence: %.2f%%\n", res.Confidence*100)
}
