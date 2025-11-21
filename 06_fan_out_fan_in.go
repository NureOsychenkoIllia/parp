package main

import (
    "fmt"
    "sync"
)

func fanOut(input <-chan int, numWorkers int) []<-chan int {
    outputs := make([]<-chan int, numWorkers)
    for i := 0; i < numWorkers; i++ {
        outputs[i] = worker(input)
    }
    return outputs
}

func worker(input <-chan int) <-chan int {
    output := make(chan int)
    go func() {
        for n := range input {
            output <- n * n // Обробка
        }
        close(output)
    }()
    return output
}

func fanIn(inputs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    output := make(chan int)

    for _, ch := range inputs {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for n := range c {
                output <- n
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(output)
    }()

    return output
}

func main() {
    input := make(chan int)
    go func() {
        for i := 1; i <= 10; i++ {
            input <- i
        }
        close(input)
    }()

    workers := fanOut(input, 3)
    results := fanIn(workers...)

    for result := range results {
        fmt.Println(result)
    }
}
