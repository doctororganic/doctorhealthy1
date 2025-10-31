package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")
	
	r.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"users":[]}`))
	}).Methods("GET")
	
	r.HandleFunc("/api/v1/foods", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"foods":[]}`))
	}).Methods("GET")
	
	r.HandleFunc("/api/v1/workouts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"workouts":[]}`))
	}).Methods("GET")
	
	r.HandleFunc("/api/v1/recipes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"recipes":[]}`))
	}).Methods("GET")
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
