// Файл: benchmark.go
// Запуск: go run benchmark.go
// Комплексний бенчмарк для порівняння послідовного та паралельного виконання

package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// ============== Тест 1: Обчислення з математичними операціями ==============
// Для кожного елемента масиву виконуємо 50 ітерацій тригонометричних функцій
// Формула: result = sin(x) * cos(x) + sqrt(|x| + 1)
// Це створює достатнє навантаження на CPU для демонстрації паралелізму

func heavyComputation(v float64) float64 {
	// Обчислювально важка операція
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

func benchmarkHeavyComputation() (time.Duration, time.Duration) {
	const size = 500_000
	arr := make([]float64, size)
	for i := range arr {
		arr[i] = rand.Float64() * 100
	}

	// Послідовно
	start := time.Now()
	_ = computeSequential(arr)
	seqTime := time.Since(start)

	// Паралельно
	start = time.Now()
	_ = computeParallel(arr, runtime.NumCPU())
	parTime := time.Since(start)

	return seqTime, parTime
}

// ============== Тест 2-3: Множення матриць ==============

func multiplySequential(a, b, c [][]float64, n int) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				c[i][j] += a[i][k] * b[k][j]
			}
		}
	}
}

func multiplyParallel(a, b, c [][]float64, n int, numWorkers int) {
	var wg sync.WaitGroup
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
				for j := 0; j < n; j++ {
					for k := 0; k < n; k++ {
						c[i][j] += a[i][k] * b[k][j]
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
			m[i][j] = rand.Float64() * 10
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

func benchmarkMatrix(size int) (time.Duration, time.Duration) {
	a := createMatrix(size)
	b := createMatrix(size)
	c1 := createZeroMatrix(size)
	c2 := createZeroMatrix(size)

	// Послідовно
	start := time.Now()
	multiplySequential(a, b, c1, size)
	seqTime := time.Since(start)

	// Паралельно
	start = time.Now()
	multiplyParallel(a, b, c2, size, runtime.NumCPU())
	parTime := time.Since(start)

	return seqTime, parTime
}

// ============== Тест 4: Worker Pool ==============

type Job struct {
	ID   int
	Data int
}

type Result struct {
	JobID  int
	Output int
}

func processJob(job Job) Result {
	// Симуляція обчислювальної роботи
	time.Sleep(100 * time.Millisecond)
	return Result{JobID: job.ID, Output: job.Data * job.Data}
}

func workerPoolSequential(jobs []Job) []Result {
	results := make([]Result, len(jobs))
	for i, job := range jobs {
		results[i] = processJob(job)
	}
	return results
}

func workerPoolParallel(jobs []Job, numWorkers int) []Result {
	jobsCh := make(chan Job, len(jobs))
	resultsCh := make(chan Result, len(jobs))
	var wg sync.WaitGroup

	// Запуск воркерів
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobsCh {
				resultsCh <- processJob(job)
			}
		}()
	}

	// Відправка завдань
	for _, job := range jobs {
		jobsCh <- job
	}
	close(jobsCh)

	// Очікування завершення
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Збір результатів
	results := make([]Result, 0, len(jobs))
	for result := range resultsCh {
		results = append(results, result)
	}
	return results
}

func benchmarkWorkerPool() (time.Duration, time.Duration) {
	const numJobs = 100
	jobs := make([]Job, numJobs)
	for i := 0; i < numJobs; i++ {
		jobs[i] = Job{ID: i, Data: i}
	}

	// Послідовно
	start := time.Now()
	_ = workerPoolSequential(jobs)
	seqTime := time.Since(start)

	// Паралельно
	start = time.Now()
	_ = workerPoolParallel(jobs, runtime.NumCPU())
	parTime := time.Since(start)

	return seqTime, parTime
}

// ============== Main ==============

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.2f мкс", float64(d.Microseconds()))
	} else if d < time.Second {
		return fmt.Sprintf("%.2f мс", float64(d.Milliseconds()))
	}
	return fmt.Sprintf("%.2f с", d.Seconds())
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║        БЕНЧМАРК: Порівняння послідовного та паралельного виконання   ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  Кількість CPU ядер: %-48d ║\n", runtime.NumCPU())
	fmt.Printf("║  GOMAXPROCS: %-56d ║\n", runtime.GOMAXPROCS(0))
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Println("┌──────────────────────────────┬────────────┬────────────┬─────────────┐")
	fmt.Println("│ Тест                         │ Послідовно │ Паралельно │ Прискорення │")
	fmt.Println("├──────────────────────────────┼────────────┼────────────┼─────────────┤")

	// Тест 1: Важкі обчислення
	fmt.Print("│ Важкі обчислення (500K)...   │")
	seqHeavy, parHeavy := benchmarkHeavyComputation()
	speedupHeavy := float64(seqHeavy) / float64(parHeavy)
	fmt.Printf(" %10s │ %10s │ %9.2fx  │\n", formatDuration(seqHeavy), formatDuration(parHeavy), speedupHeavy)

	// Тест 2: Матриці 512x512
	fmt.Print("│ Множення матриць 512x512...  │")
	seqMat512, parMat512 := benchmarkMatrix(512)
	speedupMat512 := float64(seqMat512) / float64(parMat512)
	fmt.Printf(" %10s │ %10s │ %9.2fx  │\n", formatDuration(seqMat512), formatDuration(parMat512), speedupMat512)

	// Тест 3: Матриці 1024x1024
	fmt.Print("│ Множення матриць 1024x1024...│")
	seqMat1024, parMat1024 := benchmarkMatrix(1024)
	speedupMat1024 := float64(seqMat1024) / float64(parMat1024)
	fmt.Printf(" %10s │ %10s │ %9.2fx  │\n", formatDuration(seqMat1024), formatDuration(parMat1024), speedupMat1024)

	// Тест 4: Worker Pool
	fmt.Print("│ Worker Pool (100 задач)...   │")
	seqWP, parWP := benchmarkWorkerPool()
	speedupWP := float64(seqWP) / float64(parWP)
	fmt.Printf(" %10s │ %10s │ %9.2fx  │\n", formatDuration(seqWP), formatDuration(parWP), speedupWP)

	fmt.Println("└──────────────────────────────┴────────────┴────────────┴─────────────┘")

	fmt.Println()
	fmt.Println("Висновок:")
	fmt.Printf("  • Середнє прискорення: %.2fx\n", (speedupHeavy+speedupMat512+speedupMat1024+speedupWP)/4)
	fmt.Printf("  • Теоретичний максимум (закон Амдала): ~%dx\n", runtime.NumCPU())
	fmt.Println("  • Ефективність паралелізації залежить від характеру задачі")
}
