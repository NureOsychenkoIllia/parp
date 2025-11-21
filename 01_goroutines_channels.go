package main

import (
    "fmt"
    "time"
)

func sum(arr []int, ch chan int) {
    total := 0
    for _, v := range arr {
        total += v
    }
    ch <- total
}

func main() {
    arr := make([]int, 10000000)
    for i := range arr {
        arr[i] = i
    }

    // Послідовне виконання
    start := time.Now()
    total := 0
    for _, v := range arr {
        total += v
    }
    fmt.Printf("Послідовно: %d, час: %v\n", total, time.Since(start))

    // Паралельне виконання
    start = time.Now()
    ch := make(chan int)
    mid := len(arr) / 2
    
    go sum(arr[:mid], ch)
    go sum(arr[mid:], ch)
    
    result := <-ch + <-ch
    fmt.Printf("Паралельно: %d, час: %v\n", result, time.Since(start))
}
