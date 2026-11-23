package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Store each object separately and keep only IDs in a sorted-set index. This
// avoids reading one enormous JSON array for every page.
func seed(ctx context.Context, rdb *redis.Client) error {
	pipe := rdb.Pipeline()
	for id := 1; id <= 25; id++ {
		product := Product{ID: id, Name: fmt.Sprintf("Product %02d", id)}
		value, _ := json.Marshal(product)
		pipe.Set(ctx, fmt.Sprintf("example:product:%d", id), value, 0)
		pipe.ZAdd(ctx, "example:products:by-id", redis.Z{Score: float64(id), Member: id})
	}
	_, err := pipe.Exec(ctx) // one network round trip for all writes
	return err
}

func fetchPage(ctx context.Context, rdb *redis.Client, page, pageSize int64) ([]Product, error) {
	start := (page - 1) * pageSize
	ids, err := rdb.ZRange(ctx, "example:products:by-id", start, start+pageSize-1).Result()
	if err != nil || len(ids) == 0 {
		return nil, err
	}

	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = "example:product:" + id
	}
	values, err := rdb.MGet(ctx, keys...).Result() // all page objects in one call
	if err != nil {
		return nil, err
	}

	products := make([]Product, 0, len(values))
	for _, value := range values {
		if value == nil { // object expired/deleted but its index entry remains
			continue
		}
		var product Product
		if err := json.Unmarshal([]byte(value.(string)), &product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()
	if err := seed(ctx, rdb); err != nil {
		log.Fatal(err)
	}

	products, err := fetchPage(ctx, rdb, 2, 5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("page 2: %+v\n", products)
}
