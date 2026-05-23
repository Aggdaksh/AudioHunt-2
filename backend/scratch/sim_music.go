//go:build ignore

package main

import (
	"Shazam-Vscode/backend/internal/matching"
	"fmt"
	"math/rand"
)

func main() {
	// Simulate "Shape of You" vs "Dhurandhar"
	// Generate simulated bands for Dhurandhar (repetitive bass)
	ref := make([]string, 10000)
	qry := make([]string, 300)

	// Real music has structure. Let's just assume random noise peaks at X% for 12 mismatches.
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

	bestConf := 0.0
	for offset := 0; offset <= len(ref)-len(qry); offset++ {
		matches := 0
		for i := 0; i < len(qry); i++ {
			if matching.HammingHex(ref[offset+i], qry[i]) <= 12 {
				matches++
			}
		}
		conf := float64(matches) / float64(len(qry))
		if conf > bestConf {
			bestConf = conf
		}
	}
	fmt.Printf("Random Noise (No Shift), Max Mismatches = 12, Peak = %.2f%%\n", bestConf*100)
}
