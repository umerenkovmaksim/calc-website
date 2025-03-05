package main

import (
	"calc-website/config"
	"calc-website/internal/agent"
)

func main() {
	cfg := config.LoadConfig()
	agent.Run(cfg)
}
