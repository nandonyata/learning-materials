package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Plan    string `json:"plan"`
}

var users = []User{
	{1, "Ada", "eg", "pro"}, {2, "Lin", "eg", "free"},
	{3, "Sam", "uk", "pro"}, {4, "Mona", "eg", "pro"},
}

func seed(ctx context.Context, rdb *redis.Client) error {
	pipe := rdb.Pipeline()
	for _, user := range users {
		value, _ := json.Marshal(user)
		pipe.Set(ctx, fmt.Sprintf("example:filtered:user:%d", user.ID), value, 0)
		pipe.SAdd(ctx, "example:filtered:country:"+user.Country, user.ID)
		pipe.SAdd(ctx, "example:filtered:plan:"+user.Plan, user.ID)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func queryCacheKey(filters map[string]string) string {
	parts := make([]string, 0, len(filters))
	for key, value := range filters {
		parts = append(parts, key+"="+value)
	}
	sort.Strings(parts) // map iteration is random, so canonicalize first
	hash := sha256.Sum256([]byte(strings.Join(parts, "&")))
	return "example:filtered:query:" + hex.EncodeToString(hash[:8])
}

func findUsers(ctx context.Context, rdb *redis.Client, filters map[string]string) ([]User, error) {
	cacheKey := queryCacheKey(filters)
	if value, err := rdb.Get(ctx, cacheKey).Bytes(); err == nil {
		fmt.Println("QUERY CACHE: hit", cacheKey)
		var result []User
		return result, json.Unmarshal(value, &result)
	} else if err != redis.Nil {
		return nil, err
	}
	fmt.Println("QUERY CACHE: miss", cacheKey)

	indexKeys := make([]string, 0, 2)
	if country := filters["country"]; country != "" {
		indexKeys = append(indexKeys, "example:filtered:country:"+country)
	}
	if plan := filters["plan"]; plan != "" {
		indexKeys = append(indexKeys, "example:filtered:plan:"+plan)
	}
	if len(indexKeys) == 0 {
		return nil, fmt.Errorf("at least one supported filter is required")
	}

	ids, err := rdb.SInter(ctx, indexKeys...).Result()
	if err != nil {
		return nil, err
	}
	sort.Strings(ids) // set intersections have no guaranteed order
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = "example:filtered:user:" + id
	}
	values, err := rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	result := make([]User, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}
		var user User
		if err := json.Unmarshal([]byte(value.(string)), &user); err != nil {
			return nil, err
		}
		result = append(result, user)
	}

	encoded, _ := json.Marshal(result)
	if err := rdb.Set(ctx, cacheKey, encoded, time.Minute).Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()
	if err := seed(ctx, rdb); err != nil {
		log.Fatal(err)
	}

	filters := map[string]string{"country": "eg", "plan": "pro"}
	rdb.Del(ctx, queryCacheKey(filters)) // make the lesson repeatable
	for request := 1; request <= 2; request++ {
		result, err := findUsers(ctx, rdb, filters)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("request %d: %+v\n", request, result)
	}
}
