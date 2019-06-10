package main

import (
	"math"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			burnCPU()
		}()
	}
	wg.Wait()
}

func burnCPU() {
	var v float64
	for i := 0; i < 100000000; i++ {
		v += math.Sin(float64(i))
	}
}
