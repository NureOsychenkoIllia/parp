package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const SIZE = 512

func multiplySequential(a, b, c [][]float64, n int) {
	for i := 0; i < n; i++ {
		for k := 0; k < n; k++ {
			temp := a[i][k]
			for j := 0; j < n; j++ {
				c[i][j] += temp * b[k][j]
			}
		}
	}
}

func multiplyParallel(a, b, c [][]float64, n int) {
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	rowsPerWorker := n / numWorkers

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		startRow := w * rowsPerWorker
		endRow := startRow + rowsPerWorker
		if w == numWorkers-1 {
			endRow = n
		}

		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				for k := 0; k < n; k++ {
					temp := a[i][k]
					for j := 0; j < n; j++ {
						c[i][j] += temp * b[k][j]
					}
				}
			}
		}(startRow, endRow)
	}
	wg.Wait()
}

func createMatrix(n int) [][]float64 {
	m := make([][]float64, n)
	for i := range m {
		m[i] = make([]float64, n)
		for j := range m[i] {
			m[i][j] = float64(i + j)
		}
	}
	return m
}

func main() {
	fmt.Printf("CPU cores: %d\n", runtime.NumCPU())

	a := createMatrix(SIZE)
	b := createMatrix(SIZE)
	c1 := createMatrix(SIZE)
	c2 := createMatrix(SIZE)

	start := time.Now()
	multiplySequential(a, b, c1, SIZE)
	seqTime := time.Since(start)
	fmt.Printf("Sequential: %v\n", seqTime)

	start = time.Now()
	multiplyParallel(a, b, c2, SIZE)
	parTime := time.Since(start)
	fmt.Printf("Parallel: %v\n", parTime)

	fmt.Printf("Speedup: %.2fx\n", float64(seqTime)/float64(parTime))
}
