package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// The root program is only a connection check. The real examples are under cmd/.
func main() {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal("cannot connect to Redis: ", err)
	}

	fmt.Println("Redis is ready.")
	fmt.Println("Run an example, e.g. go run ./cmd/01-cache-aside")
}
