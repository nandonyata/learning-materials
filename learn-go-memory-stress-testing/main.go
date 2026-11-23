package main

import (
	"encoding/json"
	"net/http"
	"net/http/pprof"
)

// Product represents an e-commerce product.
type Product struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Data string `json:"data"`
}

func main() {
    mux := http.NewServeMux()
    // Add pprof endpoints
    mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
    mux.Handle("/debug/pprof/heap", http.HandlerFunc(pprof.Handler("heap").ServeHTTP))
    mux.HandleFunc("/api/products", handleProducts)
    http.ListenAndServe(":8080", mux)
}

func handleProducts(w http.ResponseWriter, r *http.Request) {
    // Simulate 1MB JSON payload
    payload := make([]byte, 1024*1024)
    product := Product{
        ID:   1,
        Name: "Sample Product",
        Data: string(payload),
    }
	
    // Inefficient string concatenation
    result := ""
    for i := 0; i < 100; i++ {
        result += product.Name + " " + string(i)
    }

    resp, err := json.Marshal(product)
    if err != nil {
        http.Error(w, "Failed to marshal", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(resp)
}