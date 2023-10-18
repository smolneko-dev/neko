package main

import (
	"log"

	"github.com/smolneko-dev/neko/config"
	"github.com/smolneko-dev/neko/internal/app"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	app.Run(cfg)
}
