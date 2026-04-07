package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"service": "cache-service",
		})
	})

	http.HandleFunc("/cache/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Cache service works!",
		})
	})

	log.Println("Cache Service running on port 8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal(err)
	}
}
