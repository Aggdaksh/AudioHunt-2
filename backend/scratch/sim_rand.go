//go:build ignore

package main

import (
	"Shazam-Vscode/backend/internal/matching"
	"fmt"
	"math/rand"
)

func main() {
	// Generate 10000 random hex strings
	ref := make([]string, 10000)
	qry := make([]string, 300) // ~6 seconds

	hexChars := "0123456789abcdef"
	for i := 0; i < 10000; i++ {
		s := ""
		for j := 0; j < 16; j++ {
			s += string(hexChars[rand.Intn(16)])
		}
		ref[i] = s
	}
	for i := 0; i < 300; i++ {
		s := ""
		for j := 0; j < 16; j++ {
			s += string(hexChars[rand.Intn(16)])
		}
		qry[i] = s
	}

	maxNibbles := []int{8, 9, 10, 11, 12}
	for _, maxNibbleMismatches := range maxNibbles {
		bestConf := 0.0
		for offset := 0; offset <= len(ref)-len(qry); offset++ {
			matches := 0
			for i := 0; i < len(qry); i++ {
				// NO SHIFT
				if matching.HammingHex(ref[offset+i], qry[i]) <= maxNibbleMismatches {
					matches++
				}
			}
			conf := float64(matches) / float64(len(qry))
			if conf > bestConf {
				bestConf = conf
			}
		}
		fmt.Printf("Random Noise (No Shift), Max Mismatches = %d, Peak Confidence = %.2f%%\n", maxNibbleMismatches, bestConf*100)
	}
}
