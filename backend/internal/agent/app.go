package agent

import (
	"calc-website/config"
	"log"
)

func Run(cfg *config.Config) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	StartAgents(cfg)
	select {}
}
