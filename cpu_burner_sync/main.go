package main

import (
	"math"
)

func main() {
	for i := 0; i < 4; i++ {
		burnCPU()
	}
}

func burnCPU() {
	var v float64
	for i := 0; i < 100000000; i++ {
		v += math.Sin(float64(i))
	}
}
