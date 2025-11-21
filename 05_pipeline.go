package main

import "fmt"

func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func filter(in <-chan int, predicate func(int) bool) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            if predicate(n) {
                out <- n
            }
        }
        close(out)
    }()
    return out
}

func main() {
    // Створення конвеєра: генерація -> квадрат -> фільтр (>10)
    nums := generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    squared := square(nums)
    filtered := filter(squared, func(n int) bool { return n > 10 })

    for result := range filtered {
        fmt.Println(result)
    }
}
