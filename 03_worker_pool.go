package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job)
        time.Sleep(100 * time.Millisecond) // Симуляція роботи
        results <- job * 2
    }
}

func main() {
    const numJobs = 20
    const numWorkers = 4

    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)
    var wg sync.WaitGroup

    // Запуск воркерів
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }

    // Відправка завдань
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    // Очікування завершення
    go func() {
        wg.Wait()
        close(results)
    }()

    // Збір результатів
    for result := range results {
        fmt.Println("Result:", result)
    }
}
