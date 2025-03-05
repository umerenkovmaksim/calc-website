package main

import (
	"calc-website/config"
	"calc-website/internal/orchestrator"
)

func main() {
	cfg := config.LoadConfig()
	err := orchestrator.Run(cfg)
	if err != nil {
	}
}
