//go:build ignore

package main

import (
	"fmt"
	"github.com/madelynnblue/go-dsp/fft"
	"math"
	"math/cmplx"
	"math/rand"
)

func fft_custom(x []complex128) []complex128 {
	n := len(x)
	if n <= 1 {
		return x
	}
	even := make([]complex128, n/2)
	odd := make([]complex128, n/2)
	for i := 0; i < n/2; i++ {
		even[i] = x[2*i]
		odd[i] = x[2*i+1]
	}
	evenFFT := fft_custom(even)
	oddFFT := fft_custom(odd)
	result := make([]complex128, n)
	for i := 0; i < n/2; i++ {
		t := cmplx.Exp(complex(0, -2*math.Pi*float64(i)/float64(n))) * oddFFT[i]
		result[i] = evenFFT[i] + t
		result[i+n/2] = evenFFT[i] - t
	}
	return result
}

func main() {
	N := 2048
	input := make([]float64, N)
	for i := 0; i < N; i++ {
		input[i] = rand.Float64()
	}

	cInput := make([]complex128, N)
	for i := range input {
		cInput[i] = complex(input[i], 0)
	}

	customOut := fft_custom(cInput)
	realOut := fft.FFTReal(input)

	diff := 0.0
	for i := 0; i < N/2; i++ {
		d := cmplx.Abs(customOut[i]) - cmplx.Abs(realOut[i])
		diff += math.Abs(d)
	}
	fmt.Printf("Total Difference: %f\n", diff)
}
