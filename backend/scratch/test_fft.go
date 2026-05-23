//go:build ignore

package main

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/madelynnblue/go-dsp/fft"
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
	N := 8
	input := make([]float64, N)
	for i := range input {
		input[i] = math.Sin(float64(i))
	}

	cInput := make([]complex128, N)
	for i := range input {
		cInput[i] = complex(input[i], 0)
	}

	customOut := fft_custom(cInput)
	realOut := fft.FFTReal(input)

	fmt.Println("Custom:")
	for _, c := range customOut[:N/2+1] {
		fmt.Printf("%.2f\n", cmplx.Abs(c))
	}

	fmt.Println("\ngo-dsp:")
	for _, c := range realOut {
		fmt.Printf("%.2f\n", cmplx.Abs(c))
	}
}
