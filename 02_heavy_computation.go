package main

import (
    "fmt"
    "math"
    "runtime"
    "time"
)

// heavyComputation виконує 50 ітерацій математичних операцій
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
    arr := make([]float64, 500000)
    for i := range arr {
        arr[i] = float64(i) * 0.001
    }

    start := time.Now()
    _ = computeSequential(arr)
    fmt.Printf("Послідовно: %v\n", time.Since(start))

    start = time.Now()
    _ = computeParallel(arr, runtime.NumCPU())
    fmt.Printf("Паралельно: %v\n", time.Since(start))
}
