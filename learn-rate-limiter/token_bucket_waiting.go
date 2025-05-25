package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

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
		// fmt.Println("INI IP >> ", ip)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Lock the mutex to safely access the clients map
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(1, 3)}
		}
		clients[ip].lastSeen = time.Now()
		// limiter := clients[ip].limiter

		fmt.Println(ip, "REQS TOTAL > ", clients[ip].activeRequests)

		// NEW >>>

		// If there are already too many active requests from this IP, reject the request
		if clients[ip].activeRequests > 5 { // Allow max 5 waiting requests per user
			fmt.Println("WOI BANYAK KALI< BLOK AH")
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

		mu.Unlock()

		// Wait until a token is available
		// if err := limiter.Wait(r.Context()); err != nil {
		// 	http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		// 	return
		// }
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
