package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var database = map[int]User{1: {ID: 1, Name: "Jane Doe"}}

func loadUserFromDB(id int) (User, error) {
	fmt.Println("DATABASE: queried user", id)
	user, ok := database[id]
	if !ok {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

// getUser uses the cache-aside pattern. GET both reads the value and tells the
// backend whether the key exists: redis.Nil means a cache miss, not an outage.
func getUser(ctx context.Context, rdb *redis.Client, id int) (User, error) {
	key := fmt.Sprintf("example:cache-aside:user:%d", id)
	value, err := rdb.Get(ctx, key).Bytes()
	if err == nil {
		fmt.Println("REDIS: cache hit for", key)
		var user User
		return user, json.Unmarshal(value, &user)
	}
	if err != redis.Nil {
		return User{}, fmt.Errorf("redis GET: %w", err)
	}

	fmt.Println("REDIS: cache miss for", key)
	user, err := loadUserFromDB(id)
	if err != nil {
		return User{}, err
	}
	value, err = json.Marshal(user)
	if err != nil {
		return User{}, err
	}
	if err := rdb.Set(ctx, key, value, 5*time.Minute).Err(); err != nil {
		return User{}, fmt.Errorf("redis SET: %w", err)
	}
	return user, nil
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()
	rdb.Del(ctx, "example:cache-aside:user:1")

	for i := 1; i <= 2; i++ {
		user, err := getUser(ctx, rdb, 1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("request %d result: %+v\n\n", i, user)
	}
}
