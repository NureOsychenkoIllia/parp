// Файл: heavy_computation.go
// Запуск: go run heavy_computation.go
// Демонстрація паралельних важких обчислень

package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

// heavyComputation виконує 50 ітерацій математичних операцій
// Формула: result = sin(x) * cos(x) + sqrt(|x| + 1)
func heavyComputation(v float64) float64 {
	result := v
	for i := 0; i < 50; i++ {
		result = math.Sin(result)*math.Cos(result) + math.Sqrt(math.Abs(result)+1)
	}
	return result
}

func computeSequential(arr []float64) float64 {
	var total float64
	for _, v := range arr {
		total += heavyComputation(v)
	}
	return total
}

func computeParallel(arr []float64, numWorkers int) float64 {
	ch := make(chan float64, numWorkers)
	chunkSize := len(arr) / numWorkers

	for w := 0; w < numWorkers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if w == numWorkers-1 {
			end = len(arr)
		}

		go func(data []float64) {
			var sum float64
			for _, v := range data {
				sum += heavyComputation(v)
			}
			ch <- sum
		}(arr[start:end])
	}

	var total float64
	for i := 0; i < numWorkers; i++ {
		total += <-ch
	}
	return total
}

func main() {
	fmt.Println("=== Паралельні важкі обчислення ===")
	fmt.Printf("CPU ядер: %d\n", runtime.NumCPU())
	fmt.Println()

	const size = 500_000
	arr := make([]float64, size)
	for i := range arr {
		arr[i] = float64(i) * 0.001
	}
	fmt.Printf("Розмір масиву: %d елементів\n", size)
	fmt.Println("Операція: 50 ітерацій sin(x)*cos(x)+sqrt(|x|+1) для кожного елемента")
	fmt.Println()

	// Послідовне виконання
	fmt.Print("Послідовне обчислення... ")
	start := time.Now()
	resultSeq := computeSequential(arr)
	seqTime := time.Since(start)
	fmt.Printf("завершено за %v\n", seqTime)

	// Паралельне виконання
	fmt.Print("Паралельне обчислення... ")
	start = time.Now()
	resultPar := computeParallel(arr, runtime.NumCPU())
	parTime := time.Since(start)
	fmt.Printf("завершено за %v\n", parTime)

	// Результати
	fmt.Println()
	fmt.Println("=== Результати ===")
	fmt.Printf("Послідовний результат: %.6f\n", resultSeq)
	fmt.Printf("Паралельний результат: %.6f\n", resultPar)
	fmt.Printf("Різниця: %.10f (похибка округлення)\n", math.Abs(resultSeq-resultPar))
	fmt.Println()
	fmt.Printf("Прискорення: %.2fx\n", float64(seqTime)/float64(parTime))
	fmt.Printf("Ефективність: %.1f%%\n", float64(seqTime)/float64(parTime)/float64(runtime.NumCPU())*100)
}
