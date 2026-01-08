package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	cfg "opentelemetry/internal/config"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world this is golanghhhhh!"))
}

func main() {
	http.HandleFunc("/", HelloWorld)

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

	log.Printf("Service started on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
