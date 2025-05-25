package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/RussellLuo/slidingwindow"
	"github.com/gin-gonic/gin"
)

// seq 1 20 | xargs -I {} -P 10 curl -i http://localhost:3000/

var (
	limitPerSeq = 5
)

// RateLimiterMap stores per-IP limiters
var RateLimiterMap = struct {
	sync.RWMutex
	limiters map[string]*slidingwindow.Limiter
}{
	limiters: make(map[string]*slidingwindow.Limiter),
}

// GetLimiter returns or creates a rate limiter for an IP
func GetLimiter(ip string) *slidingwindow.Limiter {
	RateLimiterMap.Lock()
	defer RateLimiterMap.Unlock()

	if limiter, exists := RateLimiterMap.limiters[ip]; exists {
		return limiter
	}

	// Create a new rate limiter (3 requests per minute)
	limiter, _ := slidingwindow.NewLimiter(time.Minute, int64(limitPerSeq), func() (slidingwindow.Window, slidingwindow.StopFunc) {
		return slidingwindow.NewLocalWindow()
	})

	RateLimiterMap.limiters[ip] = limiter
	return limiter
}

// RateLimitMiddleware applies the limiter
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP() // Use client IP as key

		limiter := GetLimiter(ip)

		fmt.Println(">>>> IP", ip)
		fmt.Println(">>>> LIMITER", limiter)

		// Check if request is allowed
		allowed := limiter.Allow()
		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too Many Requests",
				"retry": fmt.Sprintf("%d minute", 1),
			})
			c.Abort()
			return
		}

		// Proceed if within the rate limit
		c.Next()
	}
}

func main() {
	r := gin.Default()

	// Apply the rate limit middleware
	r.Use(RateLimitMiddleware())

	// Sample API route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Request successful!"})
	})

	r.Run(":3000") // Run on port 3000
}
