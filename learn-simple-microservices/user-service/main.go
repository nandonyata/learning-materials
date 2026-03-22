// Run it by executing this command: docker compose up --build
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("GET /user", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("r.Method: %v\n", r.Method)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "hello from user-service",
		})
	})

	http.HandleFunc("GET /user/product", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("r.Method: %v\n", r.Method)

		res, err := http.Get("http://product-service:4001/product")
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"message": "error fetching product user-service-product",
				"err":     err.Error(),
			})
			return
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			json.NewEncoder(w).Encode(map[string]string{
				"message": "error  status code user-service-product",
			})
			return
		}

		rawRes, err := io.ReadAll(res.Body)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"message": "error reading body response user-service-product",
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": string(rawRes),
		})
	})

	fmt.Println("Server running...")
	http.ListenAndServe(":4000", nil)
}
