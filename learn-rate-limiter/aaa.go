package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	_ "net/http/pprof"

	"golang.org/x/time/rate"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func endpointHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	message := Message{
		Status: "Successful",
		Body:   "Hi! You've reached the API. How may I help you?",
	}
	err := json.NewEncoder(writer).Encode(&message)
	if err != nil {
		return
	}
}

type client struct {
	limiter        *rate.Limiter
	lastSeen       time.Time
	activeRequests int
}

var (
	mu      sync.Mutex
	clients = make(map[string]*client)
)

func perClientRateLimiter(next http.HandlerFunc) http.Handler {
	// Cleanup routine to remove inactive clients
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the IP address from the request
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Lock the mutex to safely access the clients map
		mu.Lock()

		// Ensure the client is initialized in the map
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(1, 3)}
		}
		clients[ip].lastSeen = time.Now()

		if ip == "::1" {
			clients[ip].activeRequests = 8

		}

		// If there are already too many active requests from this IP, reject the request
		if clients[ip].activeRequests >= 5 { // Allow max 5 waiting requests per user
			mu.Unlock() // Release lock before sending response
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		// Increment the active requests count
		clients[ip].activeRequests++

		// Decrement the active requests count when the function completes (even if there's an error)
		defer func() {
			mu.Lock()                    // Lock again to modify the activeRequests safely
			clients[ip].activeRequests-- // Decrease active request count
			mu.Unlock()
		}()

		// Unlock here after you've updated activeRequests and are about to process the request
		mu.Unlock()

		// Wait until a token is available for rate limiting
		if err := clients[ip].limiter.Wait(r.Context()); err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		// Proceed to the next handler
		next(w, r)
	})
}

func main() {
	http.Handle("/ping", perClientRateLimiter(endpointHandler))
	log.Println("Server running on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error starting server:", err)
	}
}
