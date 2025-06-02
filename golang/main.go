package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

type User struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        log.Printf("Error decoding JSON: %v", err)
        http.Error(w, "Invalid JSON input", http.StatusBadRequest)
        return
    }
    if user.Email == "" || user.Password == "" {
        log.Printf("Missing email or password")
        http.Error(w, "Email and password are required", http.StatusBadRequest)
        return
    }

    key := "user:" + user.Email

    exists, err := rdb.Exists(ctx, key).Result()
    if err != nil {
        log.Printf("Redis Exists error for key %s: %v", key, err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    if exists == 1 {
        log.Printf("User already exists: %s", user.Email)
        http.Error(w, "User already exists", http.StatusConflict)
        return
    }

    _, err = rdb.HSet(ctx, key, "email", user.Email, "password", user.Password).Result()
    if err != nil {
        log.Printf("Redis HSet error for key %s: %v", key, err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    _, _ = w.Write([]byte("Signup successful"))
}

func main() {
    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "redis:6379"
    }

    rdb = redis.NewClient(&redis.Options{
        Addr: redisAddr,
    })

    err := rdb.Ping(ctx).Err()
    if err != nil {
        log.Fatalf("Failed to connect to Redis at %s: %v", redisAddr, err)
    }
    log.Printf("Connected to Redis at %s successfully", redisAddr)

    http.HandleFunc("/signup", signupHandler)
    log.Println("Golang server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
