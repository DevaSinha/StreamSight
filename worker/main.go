package main

import (
	"log"

	"github.com/DevaSinha/StreamSight/worker/config"
	"github.com/DevaSinha/StreamSight/worker/worker"
)

func main() {
	cfg := config.LoadConfig()
	if cfg.ApiEndpoint == "" {
		log.Fatal("API_ENDPOINT environment variable must be set")
	}
	worker.RunWorker(cfg, "ws://localhost:8080/ws/alerts")

	// Keep main alive
	select {}
}
