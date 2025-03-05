package main

import (
	"calc-website/config"
	"calc-website/internal/orchestrator"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	err := orchestrator.Run(cfg)
	if err != nil {
		log.Printf("run orchestrator error: %v", err.Error())
	}
}
