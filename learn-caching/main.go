package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type User struct {
	FullName string `json:"full_name"`
	Email    string `json:"email,omitempty"`
	Password string `json:"-"`
	age      int    // unexported fields are ignored
}

func main() {
    ctx := context.Background()
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    if err := client.Set(ctx, "bike:1", "Process 134", 0).Err(); err != nil {
        panic(err)
    }

    val, err := client.Get(ctx, "bike:1").Result()
    if err != nil {
        panic(err)
    }

    fmt.Println(val)

// ----


	user := User{
		FullName: "Jane Doe",
		Email:    "",
		Password: "secret",
		age:      30,
	}

	// Convert struct to JSON bytes
	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error marshaling:", err)
		return
	}


	 if err := client.Set(ctx, "user:1", jsonData, 10*time.Minute).Err(); err != nil {
        panic(err)
    }

    valuser, err := client.Get(ctx, "user:1").Result()
    if err != nil {
        panic(err)
    }

    fmt.Println(valuser)
}
