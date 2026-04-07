package main

import (
	"encoding/json"
	"log"
	"net/http"
	
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	
	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "user-profile-service",
			"version": "1.0.0",
		})
	}).Methods("GET")
	
	// Profile endpoint
	r.HandleFunc("/api/v1/profile/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"user_id": vars["user_id"],
			"name":    "Test User",
			"status":  "active",
		})
	}).Methods("GET")
	
	log.Println("User Profile Service starting on :8083")
	log.Fatal(http.ListenAndServe(":8083", r))
}
