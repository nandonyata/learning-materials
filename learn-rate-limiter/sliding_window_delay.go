package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/RussellLuo/slidingwindow"
	"github.com/gin-gonic/gin"
)

const (
	maxDelayedRequests = 5                      // Maximum requests that can be delayed at a time
	delayDuration      = 500 * time.Millisecond // Delay for each request
)

var (
	RateLimiterMap = struct {
		sync.RWMutex
		limiters map[string]*slidingwindow.Limiter
	}{
		limiters: make(map[string]*slidingwindow.Limiter),
	}

	// Track delayed requests per IP
	delayedRequests = make(map[string]int)
	delayedMutex    sync.Mutex
)

func GetLimiter(ip string) *slidingwindow.Limiter {
	RateLimiterMap.Lock()
	defer RateLimiterMap.Unlock()

	if limiter, exists := RateLimiterMap.limiters[ip]; exists {
		return limiter
	}

	// Create a new rate limiter (10 requests per second)
	limiter, _ := slidingwindow.NewLimiter(time.Second, 10, func() (slidingwindow.Window, slidingwindow.StopFunc) {
		return slidingwindow.NewLocalWindow()
	})

	RateLimiterMap.limiters[ip] = limiter
	return limiter
}

// func RateLimitMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ip := c.ClientIP()
// 		limiter := GetLimiter(ip)

// 		delayedMutex.Lock()
// 		currentDelayed := delayedRequests[ip]
// 		if currentDelayed >= maxDelayedRequests {
// 			delayedMutex.Unlock()
// 			c.JSON(http.StatusTooManyRequests, gin.H{
// 				"error": "Too Many Requests - Max Delayed Requests Reached",
// 			})
// 			c.Abort()
// 			return
// 		}
// 		delayedRequests[ip]++
// 		delayedMutex.Unlock()

// 		// Introduce a delay before processing the request
// 		time.Sleep(delayDuration)

// 		// Check if request is allowed after delay
// 		if !limiter.Allow() {
// 			delayedMutex.Lock()
// 			delayedRequests[ip]--
// 			delayedMutex.Unlock()

// 			c.JSON(http.StatusTooManyRequests, gin.H{
// 				"error": "Too Many Requests",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		// Proceed if within the rate limit
// 		c.Next()

//			// Decrease delayed request count after processing
//			delayedMutex.Lock()
//			delayedRequests[ip]--
//			delayedMutex.Unlock()
//		}
//	}
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := GetLimiter(ip)

		delayedMutex.Lock()
		fmt.Println(">>>> ", delayedRequests[ip])
		if delayedRequests[ip] >= maxDelayedRequests {

			delayedMutex.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too Many Requests - Max Delayed Requests Reached",
			})
			c.Abort()
			return
		}
		delayedRequests[ip]++
		delayedMutex.Unlock()

		// Exponential backoff settings
		baseDelay := 50 * time.Millisecond // Initial delay
		maxDelay := 2 * time.Second        // Max delay cap
		attempts := 0

		for {
			if limiter.Allow() {
				break
			}

			// Exponential backoff (2^attempts * baseDelay)
			delayDuration := baseDelay * (1 << attempts) // 2^n * baseDelay
			if delayDuration > maxDelay {
				delayDuration = maxDelay
			}

			time.Sleep(delayDuration)
			attempts++
		}

		// Proceed with request
		c.Next()

		// Decrease delayed request count after processing
		delayedMutex.Lock()
		delayedRequests[ip]--
		delayedMutex.Unlock()
	}
}

func main() {
	r := gin.Default()
	r.Use(RateLimitMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Request successful!"})
	})

	r.Run(":3000")
}
