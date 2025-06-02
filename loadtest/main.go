package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
)

type User struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func main() {
    start := time.Now()
    var wg sync.WaitGroup
    totalRequests := 100000

    wg.Add(totalRequests)

    for i := 0; i < totalRequests; i++ {
        go func(i int) {
            defer wg.Done()
            user := User{
                Email:    fmt.Sprintf("user%d@example.com", i),
                Password: "pass123",
            }
            body, _ := json.Marshal(user)
            http.Post("http://localhost:8080/signup", "application/json", bytes.NewBuffer(body))
        }(i)
    }

    wg.Wait()
    fmt.Printf("Completed %d requests in %s\n", totalRequests, time.Since(start))
}
