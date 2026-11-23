package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// A Lua script makes INCR and EXPIRE atomic. The key automatically disappears
// after the window, so Redis does not retain every user forever.
var limitScript = redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
  redis.call("EXPIRE", KEYS[1], ARGV[1])
end
return current
`)

func allow(ctx context.Context, rdb *redis.Client, userID string, limit int) (bool, int64, error) {
	key := "example:rate-limit:user:" + userID
	count, err := limitScript.Run(ctx, rdb, []string{key}, 60).Int64()
	return count <= int64(limit), count, err
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()
	rdb.Del(ctx, "example:rate-limit:user:42")

	for request := 1; request <= 5; request++ {
		allowed, count, err := allow(ctx, rdb, "42", 3)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("request=%d count=%d allowed=%v\n", request, count, allowed)
	}
}
