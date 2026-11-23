package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	members := make([]any, 1000)
	for i := range members {
		members[i] = fmt.Sprintf("user:%d", i+1)
	}
	if err := rdb.SAdd(ctx, "example:large:user-ids", members...).Err(); err != nil {
		log.Fatal(err)
	}

	// SSCAN returns batches and a cursor. COUNT is a hint, not a guarantee.
	// Process each batch immediately instead of holding the whole set in memory.
	var cursor uint64
	total := 0
	for {
		batch, nextCursor, err := rdb.SScan(ctx, "example:large:user-ids", cursor, "user:1*", 100).Result()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("processed batch of %d matching IDs\n", len(batch))
		total += len(batch)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	fmt.Println("total matching IDs:", total)
}
