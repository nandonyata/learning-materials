// Run it by executing this command: docker compose up --build.
// Remove "--build" for production.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("GET /product", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("r.Method: %v\n", r.Method)

		// w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "hello from product-service",
		})
	})

	fmt.Println("Server running...")
	http.ListenAndServe(":4001", nil)
}
