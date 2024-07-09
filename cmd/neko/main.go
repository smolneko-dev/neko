package main

import (
	"context"
	"log"

	"github.com/smolneko-dev/neko/internal/app"
	"github.com/smolneko-dev/neko/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Printf("unable to load config: %s", err.Error())
		return
	}

	app, err := app.New(ctx, cfg)
	if err != nil {
		log.Printf("can't init neko: %s", err.Error())
		return
	}

	if err = app.Run(ctx); err != nil {
		log.Printf("neko terminated with error: %s", err.Error())
	}
}
