package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"
)

// Redis Streams are a durable event log. Unlike Pub/Sub, messages can remain
// until consumers read and acknowledge them.
func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	const stream = "example:orders"
	const group = "email-workers"
	err := rdb.XGroupCreateMkStream(ctx, stream, group, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		log.Fatal(err)
	}

	id, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		MaxLen: 1000,
		Approx: true,
		Values: map[string]any{"order_id": 123, "event": "order.created"},
	}).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("published message:", id)

	streams, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group: group, Consumer: "worker-1", Streams: []string{stream, ">"}, Count: 1,
	}).Result()
	if err != nil {
		log.Fatal(err)
	}
	for _, result := range streams {
		for _, message := range result.Messages {
			fmt.Println("processing:", message.Values)
			if err := rdb.XAck(ctx, stream, group, message.ID).Err(); err != nil {
				log.Fatal(err)
			}
		}
	}
}
