package main

import (
    "context"
    "fmt"
    "time"
)

func longOperation(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d: cancelled\n", id)
            return
        default:
            fmt.Printf("Worker %d: working...\n", id)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    for i := 1; i <= 3; i++ {
        go longOperation(ctx, i)
    }

    <-ctx.Done()
    fmt.Println("All workers cancelled due to timeout")
    time.Sleep(100 * time.Millisecond) // Дати час горутинам завершитись
}
