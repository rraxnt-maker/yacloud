package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "cache-service",
			"port":    "8082",
		})
	})

	// GET endpoint
	http.HandleFunc("/cache/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Cache service is working",
			"key":     r.URL.Path,
		})
	})

	// POST endpoint
	http.HandleFunc("/cache", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Cache SET operation successful",
			})
		}
	})

	log.Println("Cache Service starting on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
