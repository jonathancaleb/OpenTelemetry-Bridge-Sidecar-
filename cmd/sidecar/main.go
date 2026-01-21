package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	cfg "opentelemetry/internal/config"
	"opentelemetry/internal/proxy"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world this is golanghhhhh!"))
}

func main() {

	cfgPath := os.Getenv("CONFIG_FILE")
	if cfgPath == "" {
		cfgPath = "internal/config/config.yaml"
	}

	absCfgPath, err := filepath.Abs(cfgPath)
	if err != nil {
		log.Fatalf("bad config path: %v", err)
	}

	config, err := cfg.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	if err := cfg.WatchConfigFile(absCfgPath, func() {
		if newCfg, err := cfg.LoadConfig(absCfgPath); err == nil {
			config = newCfg
			log.Printf("config reloaded: %+v", config)
		} else {
			log.Printf("reload failed: %v", err)
		}
	}); err != nil {
		log.Printf("config watcher error: %v", err)
	}

	port := config.Port
	if port == "" {
		port = ":8080"
	}

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
