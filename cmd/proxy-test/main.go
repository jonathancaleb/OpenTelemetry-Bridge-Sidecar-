package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"opentelemetry/internal/proxy"
)

func main() {
	// 1. Start a simple upstream server on :3000
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[UPSTREAM] Received: %s %s", r.Method, r.URL.Path)
			fmt.Fprintf(w, "Hello from upstream! Path: %s\n", r.URL.Path)
		})
		log.Println("[UPSTREAM] Starting on :3000")
		http.ListenAndServe(":3000", nil)
	}()

	// Give upstream time to start
	time.Sleep(100 * time.Millisecond)

	// 2. Start the proxy on :8080 -> :3000
	handler, err := proxy.NewHandler(":8080", "http://localhost:3000")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[PROXY] Starting on :8080 -> :3000")
	log.Println("[TEST] Try: curl http://localhost:8080/hello")
	handler.ListenAndServe()
}
