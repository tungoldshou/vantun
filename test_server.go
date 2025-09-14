//go:build ignore

package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	// Create a simple HTTP server for testing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Simulate some processing time
		time.Sleep(10 * time.Millisecond)
		w.Write([]byte("Hello from VANTUN test server!"))
	})

	server := &http.Server{
		Addr:        ":8080",
		ReadTimeout: 5 * time.Second,
	}

	log.Println("Starting test server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
