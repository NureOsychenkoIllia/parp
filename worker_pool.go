// Файл: worker_pool.go
// Запуск: go run worker_pool.go

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	ID   int
	Data int
}

type Result struct {
	JobID  int
	Output int
	Worker int
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		processingTime := time.Duration(rand.Intn(100)) * time.Millisecond
		time.Sleep(processingTime)

		result := Result{
			JobID:  job.ID,
			Output: job.Data * job.Data,
			Worker: id,
		}
		results <- result
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const numJobs = 20
	const numWorkers = 4

	fmt.Println("=== Worker Pool Pattern ===")
	fmt.Printf("Кількість завдань: %d\n", numJobs)
	fmt.Printf("Кількість воркерів: %d\n", numWorkers)
	fmt.Println()

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)
	var wg sync.WaitGroup

	fmt.Println("Запуск воркерів...")
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	start := time.Now()
	fmt.Println("Відправка завдань...")
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Data: j}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println()
	fmt.Println("Результати:")
	for result := range results {
		fmt.Printf("  Job %2d: %d^2 = %3d (Worker %d)\n",
			result.JobID, result.JobID, result.Output, result.Worker)
	}

	elapsed := time.Since(start)
	fmt.Println()
	fmt.Printf("Загальний час: %v\n", elapsed)
	fmt.Printf("Середній час на завдання: %v\n", elapsed/numJobs)
}
