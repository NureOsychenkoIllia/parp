// Файл: matrix_multiply.go
// Запуск: go run matrix_multiply.go
// Запуск з детектором гонок: go run -race matrix_multiply.go

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const SIZE = 512

// multiplySequential виконує послідовне множення матриць
func multiplySequential(a, b, c [][]float64, n int) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				c[i][j] += a[i][k] * b[k][j]
			}
		}
	}
}

// multiplyParallel виконує паралельне множення матриць
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

func createZeroMatrix(n int) [][]float64 {
	m := make([][]float64, n)
	for i := range m {
		m[i] = make([]float64, n)
	}
	return m
}

func verifyResults(c1, c2 [][]float64, n int) bool {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if c1[i][j] != c2[i][j] {
				return false
			}
		}
	}
	return true
}

func main() {
	fmt.Println("=== Паралельне множення матриць на Go ===")
	fmt.Printf("Розмір матриці: %dx%d\n", SIZE, SIZE)
	fmt.Printf("Кількість CPU: %d\n", runtime.NumCPU())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Println()

	a := createMatrix(SIZE)
	b := createMatrix(SIZE)
	c1 := createZeroMatrix(SIZE)
	c2 := createZeroMatrix(SIZE)

	fmt.Print("Послідовне множення... ")
	start := time.Now()
	multiplySequential(a, b, c1, SIZE)
	seqTime := time.Since(start)
	fmt.Printf("завершено за %v\n", seqTime)

	fmt.Print("Паралельне множення... ")
	start = time.Now()
	multiplyParallel(a, b, c2, SIZE)
	parTime := time.Since(start)
	fmt.Printf("завершено за %v\n", parTime)

	fmt.Println()
	if verifyResults(c1, c2, SIZE) {
		fmt.Println("✓ Результати співпадають")
	} else {
		fmt.Println("✗ Результати НЕ співпадають!")
	}

	fmt.Println()
	fmt.Println("=== Статистика ===")
	fmt.Printf("Послідовний час: %v\n", seqTime)
	fmt.Printf("Паралельний час: %v\n", parTime)
	fmt.Printf("Прискорення: %.2fx\n", float64(seqTime)/float64(parTime))
	fmt.Printf("Ефективність: %.1f%%\n",
		float64(seqTime)/float64(parTime)/float64(runtime.NumCPU())*100)
}
